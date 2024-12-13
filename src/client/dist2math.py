from math import floor
from decimal import Decimal, getcontext

digits = 0

# 620 is 31 iteration
# (x^3)+25x -- 31 iterations -- 0.43  seconds
# (x^2)+20x -- 31 iterations -- 0.008 seconds
def f(x: int) -> int:
  # it will take 54 iterations to hit 0.43 seconds (0.446)
  return (x**2)+20*x

class MathFunc:
  def __init__(self) -> None:
    pass

  def SetDigitCount(newDigs: int) -> None:
    global digits
    digits = newDigs

  # S: decimal value.
  # n digits of S after o offset.
  def GetOffset(s: Decimal, o: int, n: int) -> Decimal:
    k = ((s * (Decimal(10) ** (o - 1))) - (Decimal(10) ** (o - 1)))
    return floor((Decimal(10) ** n) * (k - floor(k)))

  # same thing used in y-cruncher
  def CompSqrt2(accuracy: int) -> Decimal:
    # Fixes #2 | This is very slow, but does resolve
    # the issue. It takes 31 iterations to take 0.43 seconds
    # to compute. It does grow faster than 20x though.
    getcontext().prec = f(accuracy)
    a_n = Decimal(2)
    for _ in range(accuracy):
      a_n = (a_n / 2) + (1 / a_n)
    return a_n