# Algorithms Used In Dist2

## Algorithm 1: Newton's Method
This algorithm is used in [y-cruncher](http://www.numberworld.org/y-cruncher/) as well. It is the "defacto" method of computing square
roots. Essentially, the algorithm consists of a recursive relation:

$$
x_{n+1} = \frac{1}{2}\left(x_n + \frac{a}{x_n}\right)
$$

where $a$ is the square root we're trying to find, i.e. $a=2$ for computing
$\sqrt{2}$.

**Time Complexity**: `O(n log(n))`

## Algorithm 2: Spigot
This algorithm is used for validation. It requires less mathematical
operations, just multiplication, addition and subtraction.
It goes like this:

1. Start with a pair $\left\langle 5M, 5 \right\rangle$, where $M$ is the
number to square root.
1. Let the pair be equivalent to $\left\langle P,Q \right\rangle$.
    1. If $P\geq{Q}$, then the next sequence in the pair is $\left\langle  P-Q, Q+10\right\rangle$
    2. Otherwise, the next sequence is $\left\langle 100P,10Q-45\right\rangle$

Convergence is quadratic, however the formula below is accurate for
giving the amount of terms needed to give a certain amount of precision

$$
f(x) = \frac{(2000\cdot(x+3))-3088.28}{377}
$$

For example, $\left\lceil f(500) \right\rceil = 2661$, meaning the
algorithm needs 2,661 iterations to get at least $500$ correct digits.
Note that $f(x)$ **does not** provide a lower or upper bound, it is just
used to get a value that is in that range.


**Time Complexity**: Don't know yet.


# Sources
Algorithm 2: 
> "A Spigot-Algorithm for Square-Roots: Explained and Extended" Goldberg, Mayer. 2023 Dec 23 [doi:10.48550/arXiv.2312.15338](https://arxiv.org/abs/2312.15338)

