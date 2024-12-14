from math import floor
from decimal import Decimal, getcontext

class MathFunc:
  def __init__(self) -> None:
    pass

  # S: decimal value.
  # n digits of S after o offset.
  def GetOffset(s: Decimal, o: int, n: int) -> int:
    shifted = (s-1)*(Decimal(10)**o)
    fractional_part = shifted - floor(shifted)
    extracted_digits = floor(fractional_part*(Decimal(10)**n))
    return extracted_digits


  # same thing used in y-cruncher
  def CompSqrt2(accuracy: int) -> Decimal:
    getcontext().prec = accuracy + 10
    a_n = Decimal(2)
    for _ in range(20):
      next_a_n = (a_n / 2) + (1 / a_n)
      a_n = next_a_n
    return a_n
  
  def GetActual(off: int, digits: int) -> str:
    # Computation truncates the first digit if it's a zero.
    # so having it be an integer will remove it.
    t = str(MathFunc.GetOffset(MathFunc.CompSqrt2(off+digits), off, digits))
    
    # Fixes #3
    if len(t) == 19:
      return "0"+nn
    return nn
