# montecarlo

Calculating $\pi$ using the Monte Carlo method relies on probability and geometry. Imagine throwing thousands of darts randomly at a square board that has a circular target perfectly inscribed inside of it. By comparing the number of darts that land inside the circle to the total number thrown, we can approximate $\pi$.

## The Geometric Concept
 * **The Square:** Assume our square has side lengths of $2$, meaning its area is $4$.
 * **The Circle:** The circle inside has a radius ($r$) of $1$. The area of the circle is $\pi r^2 = \pi$.
 * **The Probability:** If you throw a dart completely at random, and it lands somewhere in the square, the 
   mathematical probability of it landing inside the circle is the ratio of their areas: $\frac{\text{Area of Circle}
   }{\text{Area of Square}} = \frac{\pi}{4}$.

## The Algorithm

Because probability dictates that the ratio of random hits will equal the ratio of the areas, we can reverse the equation to calculate $\pi$:

$$
\pi \approx 4\times \left(\frac{\text{Number of points inside the circle}}{\text{Total number of points generated}}\right)
$$

 1. **Set up boundaries:** Generate random $(x, y)$ coordinates in a 2D space, ensuring both $x$ and $y$ are numbers between $-1$ and $1$.
 2. **Count the total:** Keep track of every point generated (`Total Points`).
 3. **Check if it's inside:** Use the Pythagorean theorem $(x^2 + y^2 \leq r^2)$ to see if the point's distance from the center is less than or equal to $1$ (the radius). If it is, increment your Points Inside counter.
 4. **Calculate Pi:** After generating a massive number of points, plug your counts into the formula.

## Why does it work?

The Monte Carlo method doesn't calculate the exact value of $\pi$ right away; it gets closer as you increase the number of iterations. The accuracy scales with the square root of the number of points generated (e.g., to get three decimal places of accuracy, you will typically need around a million points).