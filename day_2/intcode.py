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


output = recursive_execute_intcode(code)
print(output)
