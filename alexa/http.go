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

// GetHTTPSkillHandler provides a http.Handler to have the freedom to use any http framework.
func (skill *Skill) GetHTTPSkillHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Validate request
		if !skill.SkipValidation {
			if err := isValidAlexaCertificate(r); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		// unmarshall body request
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body. "+err.Error(), http.StatusBadRequest)
			return
		}
		if skill.Verbose {
			log.Println("--> Request: ", string(bodyBytes))
		}

		var requestEnvelope *RequestEnvelope
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestEnvelope)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !skill.SkipValidation {
			if err := requestEnvelope.isRequestValid(skill.ApplicationID); err != nil {
				return
			}
		}

		response, err := requestEnvelope.handleRequest(skill)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error serializing response. "+err.Error(), http.StatusInternalServerError)
		}

		if skill.Verbose {
			log.Println("--> Response: ", string(json))
		}
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	})
}

func (requestEnvelope *RequestEnvelope) isRequestValid(expectedAppID string) error {
	// Check the timestamp
	if !requestEnvelope.verifyTimestamp() {
		return errors.New("Request too old to continue (>150s).")
	}

	// Check the app id
	if requestEnvelope.Context.System.Application.ApplicationID != expectedAppID {
		return errors.New("Alexa AppplicationID mismatch! Got: " + requestEnvelope.Context.System.Application.ApplicationID)
	}
	return nil
}

func isValidAlexaCertificate(r *http.Request) error {
	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		return errors.New("Invalid cert URL: " + certURL)
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		return err
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		return errors.New("Failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		return errors.New("Amazon certificate expired")
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		return errors.New("Amazon certificate invalid.")
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err = io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&bodyBuf)

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil {
		return errors.New("Signature match failed.")
	}

	return nil
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
