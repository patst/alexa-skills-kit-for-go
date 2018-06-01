package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/patst/alexa-skills-kit-for-go/alexa"
)

var breathAnimationRed = buildBreathAnimation("552200", "ff0000", 30, 1200)
var breathAnimationGreen = buildBreathAnimation("004411", "00ff00", 30, 1200)
var breathAnimationBlue = buildBreathAnimation("003366", "0000ff", 30, 1200)
var animations = [3][]alexa.GadgetAnimationStep{breathAnimationRed, breathAnimationGreen, breathAnimationBlue}

func buildBreathAnimation(fromRgbHex, toRgbHex string, steps, totalDuration int) []alexa.GadgetAnimationStep {
	halfSteps := steps / 2
	halfTotalDuration := totalDuration / 2
	return append(buildSequentialAnimation(fromRgbHex, toRgbHex, halfSteps, halfTotalDuration),
		buildSequentialAnimation(toRgbHex, fromRgbHex, halfSteps, halfTotalDuration)...)
}

/*
   The following are animation generation functions that work with the
   hexadecimal format that SetLight expects.
*/
func buildSequentialAnimation(fromRgbHex, toRgbHex string, steps, totalDuration int) []alexa.GadgetAnimationStep {
	fr, fg, fb := convert(fromRgbHex)
	tr, tg, tb := convert(toRgbHex)

	deltaRed := (tr - fr) / steps
	deltaGreen := (tg - fg) / steps
	deltaBlue := (tb - fb) / steps

	oneStepDuration := totalDuration / steps

	animationSteps := make([]alexa.GadgetAnimationStep, steps)
	for i := 0; i < steps; i++ {
		animationSteps[i] = alexa.GadgetAnimationStep{
			DurationMs: oneStepDuration,
			Color:      alexa.RgbToHex(fr, fg, fb),
			Blend:      true,
		}
		fr += deltaRed
		fg += deltaGreen
		fb += deltaBlue
	}
	return animationSteps
}

/*
   Build a 'button down' animation directive.
   The animation will overwrite the default 'button down' animation.
*/
func buildButtonDownAnimationDirective(targetGadgets []string) alexa.GadgetControllerSetLightDirective {
	return alexa.GadgetControllerSetLightDirective{
		Type:          "GadgetController.SetLight",
		Version:       1,
		TargetGadgets: targetGadgets,
		Parameters: alexa.GadgetParameters{
			Animations: []alexa.GadgetAnimation{
				{
					Repeat:       1,
					TargetLights: []string{"1"},
					Sequence: []alexa.GadgetAnimationStep{
						{
							DurationMs: 300,
							Color:      "FFFF00",
							Blend:      false,
						},
					},
				},
			},
			TriggerEvent:       "buttonDown",
			TriggerEventTimeMs: 0,
		},
	}
}

// Build a 'button up' animation directive.
func buildButtonUpAnimationDirective(targetGadgets []string) alexa.GadgetControllerSetLightDirective {
	return alexa.GadgetControllerSetLightDirective{
		Type:          "GadgetController.SetLight",
		Version:       1,
		TargetGadgets: targetGadgets,
		Parameters: alexa.GadgetParameters{
			Animations: []alexa.GadgetAnimation{
				{
					Repeat:       1,
					TargetLights: []string{"1"},
					Sequence: []alexa.GadgetAnimationStep{
						{
							DurationMs: 300,
							Color:      "00FFFF",
							Blend:      false,
						},
					},
				},
			},
			TriggerEvent:       "buttonUp",
			TriggerEventTimeMs: 0,
		},
	}
}

// Build an idle animation directive.
func buildButtonIdleAnimationDirective(targetGadgets []string, animation []alexa.GadgetAnimationStep) alexa.GadgetControllerSetLightDirective {
	return alexa.GadgetControllerSetLightDirective{
		Type:          "GadgetController.SetLight",
		Version:       1,
		TargetGadgets: targetGadgets,
		Parameters: alexa.GadgetParameters{
			Animations: []alexa.GadgetAnimation{
				{
					Repeat:       100,
					TargetLights: []string{"1"},
					Sequence:     animation,
				},
			},
			TriggerEvent:       "none",
			TriggerEventTimeMs: 0,
		},
	}
}

func convert(input string) (int, int, int) {
	var hexString string
	if strings.HasPrefix(input, "#") {
		hexString = strings.Replace(input, "#", "", 1)
	}

	if len(input) == 3 {
		hexString = fmt.Sprintf("%c%c%c%c%c%c", input[0], input[0], input[1], input[1], input[2], input[2])
	}

	d, _ := hex.DecodeString(hexString)

	return int(d[0]), int(d[1]), int(d[2])
}
