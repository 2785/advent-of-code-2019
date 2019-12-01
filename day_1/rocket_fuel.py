import math
import numpy as np

with open('./day_1/rocket_parts.txt') as f:
    parts_mass = [float(v) for v in f.readlines()]


def get_fuel_for_mass(mass):
    raw_fuel = math.floor(mass/3) - 2
    return raw_fuel if raw_fuel >= 0 else 0


parts_fuel_requirement = [get_fuel_for_mass(v) for v in parts_mass]
total_parts_fuel_requirement = np.sum(parts_fuel_requirement)

print(f"Total Fuel Requirement: {total_parts_fuel_requirement}")


def recursive_get_fuel_for_mass(remaining_mass_to_fulfill, fuel=0):
    additional_fuel_requirement = get_fuel_for_mass(remaining_mass_to_fulfill)
    if additional_fuel_requirement == 0:
        return fuel
    else:
        return recursive_get_fuel_for_mass(
            additional_fuel_requirement, fuel + additional_fuel_requirement)


parts_fuel_requirement_full = [
    recursive_get_fuel_for_mass(v) for v in parts_mass]
total_parts_fuel_requirement_full = np.sum(parts_fuel_requirement_full)

print(
    f"Total Fuel Requirement (mass of fuel taken into consideration): {total_parts_fuel_requirement_full}")
