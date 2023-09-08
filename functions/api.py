from flask import Flask, jsonify, request
import json

app = Flask(__name__)

@app.route('/.netlify/functions/api', methods=['GET', 'POST'])
def api():
    if request.method == 'GET':
        operation_code = generate_operation_code()
        return jsonify({"is_success": True, "user_id": generate_user_id(), "operation_code": operation_code})

    if request.method == 'POST':
        data = request.json
        if data is None:
            return jsonify({"is_success": False, "user_id": generate_user_id(), "message": "Invalid JSON data"})

        input_data = data.get('data')
        if not input_data or not isinstance(input_data, list):
            return jsonify({"is_success": False, "user_id": generate_user_id(), "message": "Invalid input data"})

        alphabets = [char for char in input_data if char.isalpha()]
        highest_alphabet = find_highest_alphabet(alphabets)

        response = {
            "is_success": True,
            "user_id": generate_user_id(),
            "email": "rk5532@srmist.edu.in",
            "roll_number": "RA2011030020043",
            "numbers": [char for char in input_data if char.isnumeric()],
            "alphabets": alphabets,
            "highest_alphabet": [highest_alphabet]
        }
        return jsonify(response)

def generate_user_id():
    return "rohit_kumar_27062002"

def generate_operation_code():
    return "1"

def find_highest_alphabet(alphabets):
    if not alphabets:
        return None
    return max(alphabets, key=lambda x: ord(x.lower()))
