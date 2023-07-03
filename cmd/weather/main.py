from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.firefox.options import Options
from parsing import *
import json

def get_weather_html():
    options = Options()

    options.add_argument('-headless')
    #options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3')
    driver = webdriver.Firefox(options=options)
    driver.get("https://yandex.ru/pogoda/?lat=43.11554337&lon=131.885498&via=srp")
    html_content = driver.page_source
    # save html_content to file index.html
    driver.close()

    return html_content

def save_tuple_to_json(tuple):
    dict = {'status': tuple[0], 'message': tuple[1]}
    
    with open('weather.json', 'w', encoding='utf-8') as f:
        json.dump(dict, f)

def get_test_html():
    # get html from index.html file
    with open('index.html', 'r') as f:
        html_content = f.read()
    return html_content

def get_tuple_from_json():
    with open('weather.json', 'r', encoding='utf-8') as f:
        dict = json.load(f)
    return (dict['status'], dict['message'])

def main():
    html_content = get_weather_html()
    #html_content = get_test_html() #FIXME: for test only
    weather_tuple = get_formatted_message(html_content)

if __name__ == '__main__':
    main()
    
