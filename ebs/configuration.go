/**
 *    Copyright 2018 Amazon.com, Inc. or its affiliates
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Segment ...
type Segment string

// Configuration Service segments where data can be stored in
const (
	SegmentBroadcaster Segment = "broadcaster"
	SegmentDeveloper   Segment = "developer"
	SegmentGlobal      Segment = "global"
)

// ConfigurationServiceClient is the client to communicate with the Configuration Service
type ConfigurationServiceClient struct {
	client   *http.Client
	clientID string
}

// ConfigurationParams ...
type ConfigurationParams struct {
	Segment   string `json:"segment"`
	ChannelID string `json:"channel_id,omitempty"`
	Version   string `json:"version,omitempty"`
	Content   string `json:"content"`
}

// ConfigurationResponse ...
type ConfigurationResponse map[string]ConfigurationResponseData

// ConfigurationResponseData ...
type ConfigurationResponseData struct {
	Segment struct {
		SegmentType string `json:"segment_type"`
		ChannelID   string `json:"channel_id"`
	} `json:"segment"`
	Record struct {
		Version string `json:"version"`
		Content string `json:"content"`
	} `json:"record"`
}

// SetGlobalSegment sets global configuration, available to all channels
func (c *ConfigurationServiceClient) SetGlobalSegment(animalType string) {
	config := ConfigurationParams{
		Segment: string(SegmentGlobal),
		Content: animalType,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.twitch.tv/extensions/%v/configurations/", c.clientID), b)
	if err != nil {
		fmt.Println(err)
	}

	c.setRequestHeaders(req, "")

	_, e := c.Do(req)
	if e != nil {
		log.Println(e)
	}
}

// SetDeveloperSegment sets channel-specfic configuration
func (c *ConfigurationServiceClient) SetDeveloperSegment(channelID, animalFact string) {
	config := ConfigurationParams{
		Segment:   string(SegmentDeveloper),
		ChannelID: channelID,
		Content:   animalFact,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.twitch.tv/extensions/%v/configurations/", c.clientID), b)
	if err != nil {
		fmt.Println(err)
	}

	c.setRequestHeaders(req, channelID)

	_, e := c.Do(req)
	if e != nil {
		log.Println(e)
	}
}

// GetBroadcasterSegment retrieves channel-specific configuration
func (c *ConfigurationServiceClient) GetBroadcasterSegment(channelID string) (animalType string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/extensions/%v/configurations/channels/%v", c.clientID, channelID), nil)
	if err != nil {
		log.Println(err)
	}

	c.setRequestHeaders(req, channelID)

	body, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}

	var config ConfigurationResponse
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Println(err)
	}

	for _, v := range config {
		if v.Segment.SegmentType == string(SegmentBroadcaster) {
			animalType = v.Record.Content
		}
	}

	return animalType
}

// Do the actual request
func (c *ConfigurationServiceClient) Do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		currentTime := time.Now().Unix()
		resetTime, err := strconv.Atoi(resp.Header.Get("Ratelimit-Reset"))
		if err != nil {
			return nil, err
		}

		if currentTime < int64(resetTime) {
			timeDiff := time.Duration(int64(resetTime) - currentTime)
			if timeDiff > 0 {
				log.Println("Waiting on rate limit to pass")
				time.Sleep(timeDiff * time.Second)
				c.Do(req)
			}
		}
	}

	return nil, nil
}

func (c *ConfigurationServiceClient) setRequestHeaders(req *http.Request, channelID string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", NewJWT(channelID)))
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Content-Type", "application/json")
}
