import math

test_input_1 = "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51"
test_input_2 = "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7"


def get_whole_path(path_text, initial_pos=(0, 0)):
    steps_in_text = [v.strip() for v in path_text.split(",")]
    parsed_steps = [(v[0], int(v[1:])) for v in steps_in_text]
    all_points_on_path = [initial_pos]
    for step in parsed_steps:
        curr_pos = all_points_on_path[-1]
        step_range = range(1, step[1] + 1)
        if step[0] == 'U':
            all_points_on_path.extend(
                [(curr_pos[0], curr_pos[1] + v) for v in step_range])
        elif step[0] == "D":
            all_points_on_path.extend(
                [(curr_pos[0], curr_pos[1] - v) for v in step_range])
        elif step[0] == "R":
            all_points_on_path.extend(
                [(curr_pos[0] + v, curr_pos[1]) for v in step_range])
        elif step[0] == "L":
            all_points_on_path.extend(
                [(curr_pos[0] - v, curr_pos[1]) for v in step_range])
        else:
            return {"error": True, "message": f"Unknown command {step[0]} encountered"}, all_points_on_path

    return {"error": False, "message": ''}, all_points_on_path


def find_next_intersection_of_two_wires(wire_1_path, wire_2_path):
    e1, path_1 = get_whole_path(wire_1_path)
    e2, path_2 = get_whole_path(wire_2_path)
    if e1["error"]:
        return e1["message"]
    if e2["error"]:
        return e2["message"]
    intersections = set(path_1).intersection(path_2)
    sorted_intersections = sorted(
        [{"point": v, "distance": abs(v[0]) + abs(v[1])} for v in intersections], key=lambda p: p["distance"])
    if len(sorted_intersections) > 1:
        return f"Next intersection happens at {sorted_intersections[1]['point']} with distance {sorted_intersections[1]['distance']}"
    else:
        return "No other intersection found"


# Test Case: Expect 135
# test_case = find_next_intersection_of_two_wires(test_input_1, test_input_2)
# print(test_case)

# Actual Input:
with open("./day_3/inputs.txt") as f:
    lines = f.readlines()
    if len(lines) != 2:
        print("more than 2 lines found in file, abort")
    else:
        print(find_next_intersection_of_two_wires(*lines))
