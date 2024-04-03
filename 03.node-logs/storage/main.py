from flask import Flask, request, jsonify
import requests
import boto3
from botocore.exceptions import NoCredentialsError

app = Flask(__name__)

# MinIO Configuration
minio_client = boto3.client('s3',
                            endpoint_url='http://localhost:9000', # Change to your MinIO endpoint
                            aws_access_key_id='yourMinioAccessKey', # Change to your MinIO access key
                            aws_secret_access_key='yourMinioSecretKey', # Change to your MinIO secret key
                            region_name='us-east-1')

bucket_name = 'weather-data'

def store_data_in_minio(bucket, object_name, data):
    try:
        minio_client.put_object(Bucket=bucket, Key=object_name, Body=data)
        return True
    except NoCredentialsError:
        return False

@app.route('/address/<address>', methods=['GET'])
def get_weather(address):
    api_url = f"http://api.weatherapi.com/v1/current.json?key=5a2d6a9bcdd54cdd97a153506242003&q={address}"
    response = requests.get(api_url)
    if response.status_code == 200:
        # Assuming you want to store the JSON response as a string in MinIO
        data = response.text
        object_name = f"weather_{address}.json"
        if store_data_in_minio(bucket_name, object_name, data):
            return jsonify({"message": "Data stored successfully"}), 200
        else:
            return jsonify({"error": "Failed to store data in MinIO"}), 500
    else:
        return jsonify({"error": "Failed to fetch weather data"}), response.status_code

if __name__ == '__main__':
    app.run(debug=True)
