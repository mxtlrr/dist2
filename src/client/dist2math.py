from math import floor
from decimal import Decimal, getcontext

class MathFunc:
  def __init__(self) -> None:
    pass

  # S: decimal value.
  # n digits of S after o offset.
  def GetOffset(s: Decimal, o: int, n: int) -> Decimal:
    k = ((s * (Decimal(10) ** (o - 1))) - (Decimal(10) ** (o - 1)))
    return floor((Decimal(10) ** n) * (k - floor(k)))

  # same thing used in y-cruncher
  # def CompSqrt2(accuracy: int) -> float:
  #   # a_n/2 + 1/a_n
  #   a_n = 2
  #   for i in range(accuracy):
  #     a_n = a_n/2 + 1/a_n
  #   return a_n
  def CompSqrt2(accuracy: int) -> Decimal:
    getcontext().prec = accuracy + 2  # Set precision higher than desired
    a_n = Decimal(2)
    for _ in range(accuracy):
      a_n = (a_n / 2) + (1 / a_n)
    return a_n