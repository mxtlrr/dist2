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

class ValidationF:
    def __init__(self) -> None:
        pass

    # https://arxiv.org/pdf/2312.15338
    # Adapted for M=2
    def Spigot(acc: Decimal) -> str:
        pair = (Decimal(10), Decimal(5))
        for _ in range(acc):
            p=pair[0];q=pair[1]
            pair = (p-q,q+10) if p>=q else (100*p, (10*q)-45)
        return str(pair[1])[1:]

    # Produces a number that is known to calculate correct digits.
    def Accuracy(dig: int) -> Decimal:
        return Decimal((((2000*(dig+3)) - 3088.28)/377))

    def Offset(s: str, o: int, n: int) -> str:
        return s[o:n]

    def Check(newton: str, off: int, n=3) -> bool:
        digits = len(newton)
        z = ValidationF.Offset(ValidationF.Spigot(ceil(
                ValidationF.Accuracy(digits+5))), off, digits)
        last = z[len(z)-n:]
        if newton[digits-n:] == last:
            return True
        return False
