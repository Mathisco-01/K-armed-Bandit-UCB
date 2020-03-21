package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/wcharczuk/go-chart"
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
	nextStep      int
	totalScore    float64 // required for calculating average reward
	exploration   float64
}

// Returns GaussianDistribution object
func CreateGaussianDistribution(mu, sigma float64) GaussianDistribution {
	return GaussianDistribution{mu: mu, sigma: sigma}
}

// Returns random number from GaussianDistribution object
func RandomGaussian(distr GaussianDistribution) float64 {
	u1 := RandomFloat64Range(0, 1)
	u2 := RandomFloat64Range(0, 1)

	z1 := (math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2))
	return z1 * distr.sigma * distr.mu
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
	var maxi int
	for i, num := range a {
		if num > max {
			max = num
			maxi = i
		}
	}
	return maxi
}

// Returns agent object with distributions instantiated
func CreateAgent(k int) Agent {
	var agent Agent
	for i := 0; i < k; i++ {
		agent.distributions = append(agent.distributions, CreateGaussianDistribution(RandomFloat64Range(-5, 5), RandomFloat64Range(-3, 3)))
	}
	agent.q = make([]float64, len(agent.distributions))
	agent.numExplored = make([]float64, len(agent.distributions))
	agent.nextStep = rand.Intn(len(agent.distributions))
	return agent
}

func ClearAgent(agent Agent) Agent {
	agent.q = make([]float64, len(agent.distributions))
	agent.numExplored = make([]float64, len(agent.distributions))
	agent.nextStep = rand.Intn(len(agent.distributions))
	agent.totalScore = 0

	return agent
}

func AgentStep(agent Agent) Agent {

	reward := RandomGaussian(agent.distributions[agent.nextStep])

	agent.numExplored[agent.nextStep] += float64(1)
	agent.totalScore += reward
	agent.q[agent.nextStep] += (reward - agent.q[agent.nextStep]) / agent.numExplored[agent.nextStep]

	if ArraySumFloat64(agent.numExplored) > math.Inf(-1) { // do greedy search on step

		var expectedQ []float64

		for i := 0; i < len(agent.distributions); i++ {
			expectedQ = append(expectedQ, agent.q[i]+(float64(agent.exploration)*(math.Sqrt((math.Log(math.E)*ArraySumFloat64(agent.numExplored))/agent.numExplored[i]))))
		}
		agent.nextStep = RandomArgmax(expectedQ)
	}

	return agent
}

func main() {
	rand.Seed(time.Now().UnixNano())

	distributions := 50
	iterations := 200

	agent := CreateAgent(distributions)

	var explorationArray []float64
	for e:=1;e<200;e++ {
		explorationArray = append(explorationArray, (0.0 + 0.001 * float64(e)))
	}

	for _, exploration := range explorationArray {
		var averages []float64
		var xAxis []float64
		agent.exploration = exploration
		for i := 0; i < iterations; i++ {
			agent = AgentStep(agent)
			averages = append(averages, agent.totalScore/ArraySumFloat64(agent.numExplored))
			xAxis = append(xAxis, float64(i))
		}
		graph := chart.Chart{
			Title: fmt.Sprintf("Exploration: %f", agent.exploration),
			YAxis: chart.YAxis{
				Name: "Cumulative average reward",
				Range: &chart.ContinuousRange{
					Min: -2,
					Max: 5,
				},
			},
			XAxis: chart.XAxis{
				Name: "Iteration",
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: xAxis,
					YValues: averages,
				},
			},
		}


		filename := fmt.Sprintf("imgs/explr%f.png", agent.exploration)
		f, _ := os.Create(filename)
		defer f.Close()
		err := graph.Render(chart.PNG, f)
		if err != nil {
			fmt.Println(err)
		}

		agent = ClearAgent(agent)
	}

}
