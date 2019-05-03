package main

import (
	"encoding/json"
	"io"
        "io/ioutil"
	"fmt"
	"net/http"
	"strconv"
	"time"
        "github.com/nlopes/slack"
)

type InfectedFile struct {
	Name string `json:"name"` //file name
	Path string `json:"path"`  // artifact path in Artifactory
	SHA256 string `json:"sha256"` // artifact SHA 256 checksum
	Depth int `json:"depth"`  // Artifact depth in its hierarchy
	ParentSHA string `json:"parent_sha"` // Parent artifact SHA1 checksum
	DisplayName string `json:"display_name"`
	PackageType string `json:"pkg_type"`
}

type ImpactedArtifact struct {
	Name string `json:"name"` //artifact name
	DisplayName string `json:"display_name"` //issue type Artifact display name
	Path string `json:"path"`  // artifact path in Artifactory
	PackageType string `json:"pkg_type"`
	SHA256 string `json:"sha256"` // artifact SHA 256 checksum
	SHA1 string `json:"sha1"`
	Depth int `json:"depth"`  // Artifact depth in its hierarchy
	ParentSHA string `json:"parent_sha"` // Parent artifact SHA1 checksum
	InfectedFiles InfectedFiles `json:"infected_files"`
}

type Issue struct {
	Severity string `json:"severity"`
	Type string `json:"type"` //issue type license/security
	Summary string `json:"summary"`
	Description string `json:"description"`
	ImpactedArtifacts ImpactedArtifacts `json:"impacted_artifacts"`
}

type Violation struct {
  Created string `json:"created"`
  TopSeverity string `json:"top_severity"`
  WatchName string `json:"watch_name"`
  PolicyName string `json:"policy_name"`
  Issues Issues `json:"issues"`
}

type InfectedFiles []InfectedFile
type ImpactedArtifacts []ImpactedArtifact
type Issues []Issue

func SendSlack(w http.ResponseWriter, r *http.Request) {
	var violation Violation

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 5048576))
        if err != nil{
                panic(err)
        }

        if err := json.Unmarshal(body, &violation); err != nil {
                w.Header().Set("Content-Type", "application/json; charset=UTF-8")
                w.WriteHeader(200) // unprocessable entity
                if err := json.NewEncoder(w).Encode(err); err != nil {
                        panic(err)
                }
        }

	fmt.Println("Sending a Slack Message")
	//fmt.Println(len(violation.Issues[0]))
	//for i := 0; i < len(violation.Issues); i++ {
	//	for j := 0; j < len(violation.Issues[].ImpactedArtifacts); j++ {
	//		sum += i
//		}
//	}

    attachment := slack.Attachment{
      Color:         "good",
      Fallback:      "You successfully posted by Incoming Webhook URL!",
      Text:          "Violation was found with top severity: " + violation.TopSeverity + " in file: " + violation.Issues[0].ImpactedArtifacts[0].InfectedFiles[0].Name,
      Footer:        "slack api",
      FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
      Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
    }
    msg := slack.WebhookMessage{
      Attachments: []slack.Attachment{attachment},
    }

    slack.PostWebhook("https://hooks.slack.com/services/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", &msg)

}

