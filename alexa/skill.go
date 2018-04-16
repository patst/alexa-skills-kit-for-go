package alexa

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Skill configures the different Handlers for skill execution.
type Skill struct {
	ApplicationID      string
	OnLaunch           func(*LaunchRequest, *OutgoingResponse)
	OnIntent           func(*IntentRequest, *OutgoingResponse)
	OnSessionEnded     func(*SessionEndedRequest, *OutgoingResponse)
	OnAudioPlayerState func(*AudioPlayerRequest, *OutgoingResponse)
	OnSystemException  func(*SystemExceptionEncounteredRequest, *OutgoingResponse)
	OnGameEngineEvent  func(*GameEngineInputHandlerEventRequest, *OutgoingResponse)
}

// GetSkillHandler provides a http.Handler to have the freedom to use any http framework.
func (skill *Skill) GetSkillHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isDev := string(r.URL.Query().Get("dev")) == "true"
		//Validate request
		if !isValidAlexaCertificate(w, r, isDev) {
			return
		}

		// unmarshall body request
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			HTTPError(w, err.Error(), "Error reading request body. Bad Request", 400)
			return
		}
		log.Println("--> Request: ", string(bodyBytes))

		var requestEnvelope *RequestEnvelope
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestEnvelope)

		if err != nil {
			HTTPError(w, err.Error(), "Bad Request", 400)
			return
		}

		if !isRequestValid(requestEnvelope, skill.ApplicationID, isDev, w) {
			return
		}

		response, err := handleRequest(requestEnvelope, skill)

		if err != nil {
			HTTPError(w, err.Error(), "Bad Request", 400)
			return
		}

		json, err := json.Marshal(response)
		if err != nil {
			log.Fatal("Error serializing response", err)
			http.Error(w, "unexpected exception.", http.StatusInternalServerError)
		}

		log.Println("--> Response: ", string(json))
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	})
}

func isRequestValid(requestEnvelope *RequestEnvelope, expectedAppID string, isDev bool, w http.ResponseWriter) bool {
	// Check the timestamp
	if !requestEnvelope.VerifyTimestamp() && !isDev {
		HTTPError(w, "Request too old to continue (>30s).", "Bad Request", 400)
		return false
	}

	// Check the app id
	if requestEnvelope.Context.System.Application.ApplicationID != expectedAppID {
		HTTPError(w, "Alexa AppplicationId mismatch!", "Bad Request", 400)
		return false
	}
	return true
}

func handleRequest(requestEnvelope *RequestEnvelope, skill *Skill) (*OutgoingResponse, error) {
	//Read the type for this request to do the correct routing
	var commonRequest CommonRequest
	err := requestEnvelope.GetTypedRequest(&commonRequest)
	if err != nil {
		return nil, err
	}

	requestType := commonRequest.Type

	// Create response and map the session attributes from the request
	response := NewOutgoingResponse(requestEnvelope.Session.Attributes)

	// Request handling
	if requestType == "LaunchRequest" {

		if skill.OnLaunch != nil {
			//Map to the correct type
			var request LaunchRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnLaunch(&request, response)
		}
	} else if requestType == "IntentRequest" {
		if skill.OnIntent != nil {
			var request IntentRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnIntent(&request, response)
		}
	} else if requestType == "SessionEndedRequest" {
		if skill.OnSessionEnded != nil {
			var request SessionEndedRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnSessionEnded(&request, response)
		}
	} else if strings.HasPrefix(requestType, "AudioPlayer.") {
		if skill.OnAudioPlayerState != nil {
			var request AudioPlayerRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnAudioPlayerState(&request, response)
		}
	} else if strings.HasPrefix(requestType, "GameEngine.") {
		if skill.OnGameEngineEvent != nil {
			var request GameEngineInputHandlerEventRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnGameEngineEvent(&request, response)
		}
	} else if requestType == "System.ExceptionEncountered" {
		if skill.OnSystemException != nil {
			var request SystemExceptionEncounteredRequest
			// Create concrete types
			requestEnvelope.GetTypedRequest(&request)
			skill.OnSystemException(&request, response)
		}
	} else {
		return nil, errors.New("Invalid request")
	}
	return response, nil
}

// HTTPError logs the logMsg and sets the given errCode
func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Println(logMsg)
	}

	http.Error(w, err, errCode)
}

func isValidAlexaCertificate(w http.ResponseWriter, r *http.Request, isDev bool) bool {
	if isDev {
		return true
	}
	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", 401)
		return false
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
		return false
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		HTTPError(w, "Amazon certificate invalid.", "Not Authorized", 401)
		return false
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err = io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
	if err != nil {
		HTTPError(w, err.Error(), "Internal Error", 500)
		return false
	}
	//log.Println(bodyBuf.String())
	r.Body = ioutil.NopCloser(&bodyBuf)

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil {
		HTTPError(w, "Signature match failed.", "Not Authorized", 401)
		return false
	}

	return true
}

func readCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("Could not download Amazon cert file")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("Could not read Amazon cert file")
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}
