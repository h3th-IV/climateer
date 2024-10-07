import cdsapi
import sys
import json
import cfgrib
import xarray

def convert_grib_to_csv(grib_file, output_file):
    dataset = xarray.open_dataset(grib_file, engine='cfgrib')
    dataframe = dataset.to_dataframe()
    dataframe.to_csv(output_file, index=False)

def fetch_era5_reanalysis(variable, year, month, day, time, area):
    client = cdsapi.Client()
    
    # Convert inputs from strings to lists
    variable_list = variable.split(',')
    year_list = year.split(',')
    month_list = month.split(',')
    day_list = day.split(',')
    time_list = time.split(',')
    area_list = list(map(float, area.split(',')))
    
    dataset = "reanalysis-era5-single-levels"
    file_name = f"era5_{variable}_{year}_{month}_{day}_{time}.grib"
    file_csv = f"era5_{variable}_{year}_{month}_{day}_{time}.csv"
    request = {
        "product_type": ["reanalysis"],
        "variable": variable_list,
        "year": year_list,
        "month": month_list,
        "day": day_list,
        "time": time_list,
        "data_format": "grib",
        "download_format": "unarchived",
        "area": area_list
    }
    
    client.retrieve(dataset, request).download(file_name)
    convert_grib_to_csv(file_name, file_csv)
    return file_csv

if __name__ == "__main__":
    variable = sys.argv[1]
    year = sys.argv[2]
    month = sys.argv[3]
    day = sys.argv[4]
    time = sys.argv[5]
    area = sys.argv[6]

    try:
        file_name = fetch_era5_reanalysis(variable, year, month, day, time, area)
        response = {
            "message": "Data fetched successfully",
            "file_name": file_name
        }
        print(json.dumps(response))
    except Exception as e:
        error_response = {
            "error": str(e),
            "message": "Failed to fetch data"
        }
        print(json.dumps(error_response))
#to run: python reanalysis.py 2m_temperature 2023 06 01 00:00 "4,-20,20,15"