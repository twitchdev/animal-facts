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
	"math/rand"
	"time"
)

// FactsServiceClient is a local service that handles requests for animal facts.
type FactsServiceClient struct{}

// AnimalType is an enum that corresponds to animal types.
type AnimalType string

// AnimalType values
const (
	Cat AnimalType = "cat"
	Dog AnimalType = "dog"
)

// GetRandomFact returns a random cat or dog fact.
func (c *FactsServiceClient) GetRandomFact(t AnimalType) string {
	var fact string

	// Taken from https://catfact.ninja/
	catFacts := []string{
		"In the 1750s, Europeans introduced cats into the Americas to control pests.",
		"Statistics indicate that animal lovers in recent years have shown a preference for cats over dogs!",
		"Like humans, cats tend to favor one paw over another",
		"A female cat is called a queen or a molly.",
		"It has been scientifically proven that stroking a cat can lower one's blood pressure.",
	}

	// Taken from https://github.com/kinduff/dog-api/blob/master/db/seeds.rb
	dogFacts := []string{
		"Dogs have sweat glands in between their paws.",
		"During the Middle Ages, Great Danes and Mastiffs were sometimes suited with armor and spiked collars to enter a battle or to defend supply caravans.",
		"Your pup reaches his full size between 12 and 24 months.",
		"Chihuahuas are born with soft spots in their skulls, just like human babies.",
		"Obesity is the top health problem among dogs.",
	}

	rand.Seed(time.Now().Unix())

	switch t {
	case Cat:
		fact = catFacts[rand.Intn(len(catFacts))]
	case Dog:
		fact = dogFacts[rand.Intn(len(dogFacts))]
	}

	return fact
}

// GetDefaultFact returns a fact that is applicable to all animal types.
func (c *FactsServiceClient) GetDefaultFact() string {
	return "Animals are our friends."
}
