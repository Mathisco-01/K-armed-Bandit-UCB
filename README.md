# K-armed-Bandit-UCB
K-armed Bandit with Upper-Confidence-Bound action selection based on this formula.
![formula](https://github.com/Mathisco-01/K-armed-Bandit-UCB/blob/master/imgs/formula.png?raw=true)
Where *Qt(a)* is the quality of action *a* at timestep *t* (or expected reward based on the cumulative average of previous results), *c* is the **exploration** constant (where a higher *c* means a larger bias towards exploring rather than exploiting) and where *Nt(a)* is the amount of times action *a* has been picked. The last part of the algorithm incentivises the exploration of actions that haven't been explored that much.

![gif](https://github.com/Mathisco-01/K-armed-Bandit-UCB/blob/master/exploration.gif?raw=true)
Here is a fun gif demonstrating what changing the c (exploration) parameter does.
