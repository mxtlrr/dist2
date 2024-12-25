from math import floor, ceil
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
      return "0"+t
    return t

NEW_OFFSET = 50
class ValidationF:
  def __init__(self) -> None:
    pass

  # This returns a number that DOES provide the amount of correct digits,
  # not necessarily being the lower bound, but also not necessarily being
  # the upper bound.
  def CalcAccuracy(dig: int) -> Decimal:
    return Decimal( (((2000*(dig+3)) - 3088.28)/377))
    #return (((2000*(dig+3)) - 3088.28)/377)

  # Correct digits for X accuracy:
  # f(x) = 0.344262x^(0.878332) if x > 10
  # Source: https://arxiv.org/pdf/2312.15338
  def R(acc: Decimal) -> str:
    getcontext().prec = acc + 30 # Update accuracy
    pair = (Decimal(10), Decimal(5))
    for _ in range(acc):
        p=pair[0];q=pair[1]
        pair = (p-q,q+10) if p>=q else (100*p, (10*q)-45)
    return str(pair[1])[1:]

  def New(s,o,n):
    return s[o:o+n]

  # offset is the start of the newton string. n is the
  # amount of digits to check.
  def Validate(newton: str, off: int, n=3) -> tuple:
    digits = len(newton)
    kc = ValidationF.R(ceil(ValidationF.CalcAccuracy(digits+off+NEW_OFFSET)))
    z = ValidationF.New(kc, off, digits)
    last = z[len(z)-n:]
    if newton[digits-n:] == last:
      return ("", True)
    return (last, False)
