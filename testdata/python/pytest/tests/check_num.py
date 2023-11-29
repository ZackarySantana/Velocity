from src.num import sum
import pytest

def test_check_number():
    assert sum(2, 3) == 5
    assert sum(3, 3) == 6