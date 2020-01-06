package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type GaussianDistribution struct {
	mu    float64 // mean
	sigma float64 // standard deviation
}

type Agent struct {
	distributions []GaussianDistribution // holds GaussianDistribution objects
	q             []float64              // holds the Q variables for their respective GaussianDistributions
	numExplored   []float64              // times the distributions has been expored/exploited, is float64 for ease of use
	totalScore		float64								// required for calculating average reward
}

// Returns GaussianDistribution object
func CreateGaussianDistribution(mu, sigma float64) GaussianDistribution {
	return GaussianDistribution{mu: mu, sigma: sigma}
}

// Returns random number from GaussianDistribution object
func RandomGaussian(distr GaussianDistribution) float64 {
	return distr.mu + distr.sigma*RandomFloat64Range(-1, 1)
}

// Returns Random float64 within a range of (min, max)
func RandomFloat64Range(min, max int) float64 {
	return float64(min) + rand.Float64()*(float64(max)-float64(min))
}

// Returns the sum of an []int
func ArraySum(a []int) (sum int) {
	for _, num := range a {
		sum += num
	}
	return sum
}

// Returns the float64 sum of an []float64
func ArraySumFloat64(a []float64) (sum float64) {
	for _, num := range a {
		sum += num
	}
	return sum
}

// Returns the max index of an array
// In the case of identical max's it randomly selects one
func RandomArgmax(a []float64) int {

	max := math.Inf(-1)

	var maxList []int
	for i, num := range a {
		if num > max {
			max = num
			maxList = make([]int, i)
		} else if num == max {
			maxList = append(maxList, i)
		}
	}

	if len(maxList) > 0 {
		return maxList[rand.Intn(len(maxList))]
	} else {
		return 0
	}
}

// Returns agent object with distributions instantiated
func CreateAgent(k int) Agent {
	var agent Agent
	for i := 0; i < k; i++ {
		agent.distributions = append(agent.distributions, CreateGaussianDistribution(RandomFloat64Range(0, 5), RandomFloat64Range(0, 5)))
		agent.q = append(agent.q, 0)
		agent.numExplored = append(agent.numExplored, 0)
	}

	return agent
}

func AgentStep(agent Agent) Agent {
	var nextStep int

	if ArraySumFloat64(agent.numExplored) > 0 { // do greedy search on step
		exploration := 2
		var expectedQ []float64

		for i := 0; i < len(agent.distributions); i++ {
			expectedQ = append(expectedQ, agent.q[i]+float64(exploration)*(math.Sqrt(math.Log(math.E)/agent.numExplored[i])))

		}
		fmt.Println(expectedQ)
		nextStep = RandomArgmax(expectedQ)

	} else {
		nextStep = RandomArgmax(agent.q)
	}

	reward := RandomGaussian(agent.distributions[nextStep])

	agent.numExplored[nextStep] += float64(1)
	agent.totalScore += reward

	agent.q[nextStep] = reward - agent.q[nextStep] / agent.numExplored[nextStep]

	return agent
}

func main() {
	rand.Seed(time.Now().UnixNano())
	agent := CreateAgent(10)
	fmt.Printf("Number of distributions: %v\n", len(agent.distributions))
	for i := 0; i < 100000; i++ {
		fmt.Println(i)
		agent = AgentStep(agent)
	}
	fmt.Println(agent.totalScore)
	fmt.Println(agent.totalScore / ArraySumFloat64(agent.numExplored))
}
