import copy
from scipy.optimize import fsolve

with open("./day_2/input.txt") as f:
    line = f.readline()

code = [int(v) for v in line.split(',', -1)]
print(code)


def recursive_execute_intcode(inputs, curr_pos=0):
    l = len(inputs) - 1
    if curr_pos > l:
        return {"error": True, "message": "Current Position exceeds input length"}, inputs

    if inputs[curr_pos] == 99:
        return {"error": False, "message": ""}, inputs

    if inputs[curr_pos] != 1 and inputs[curr_pos] != 2:
        return {"error": True, "message": f"Current position {curr_pos} is not opcode"}, inputs

    if curr_pos + 3 > l:
        return {"error": True, "message": "Operands exceed input length"}, inputs

    operation = (lambda x, y: x +
                 y) if inputs[curr_pos] == 1 else (lambda x, y: x*y)

    first = inputs[curr_pos + 1]
    second = inputs[curr_pos + 2]
    destination = inputs[curr_pos + 3]

    if first > l or second > l or destination > l:
        return {"error": True, "message": "Operands point at position outside range"}, inputs

    inputs[destination] = operation(inputs[first], inputs[second])
    return recursive_execute_intcode(inputs, curr_pos + 4)


def get_intcode_output(in1, in2, arr):
    arr[1] = in1
    arr[2] = in2
    return recursive_execute_intcode(arr)[1][0]


output = recursive_execute_intcode(copy.deepcopy(code))
print(code)
print(output)

all_inputs = [(i, j) for i in range(100) for j in range(100)]
# print(all_inputs)
all_results = [{"input_pair": v, "output": get_intcode_output(
    v[0], v[1], copy.deepcopy(code))} for v in all_inputs]
matching_results = [v for v in all_results if v["output"] == 19690720]
final_out = 'None' if len(
    matching_results) == 0 else matching_results[0]['input_pair'][0] * 100 + matching_results[0]['input_pair'][1]
print(
    f"Part b result: {final_out}")
