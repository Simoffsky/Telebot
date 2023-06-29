from bs4 import BeautifulSoup
import requests


def get_formatted_message(html):
    try: 
        message = 'Погода во Владивостоке: \n'
        soup = BeautifulSoup(html, "html.parser")
        if not check_valid(soup): #FIXME:
            return (False, "Не удалось распарсить сайт с погодой :(. Надо чинить парсер")
        

        message += soup.find('div', class_='link__condition day-anchor i-bem').text + '\n'
        
        time = soup.find('time', class_='time fact__time').text
        message += time + '\n'
        fact_temp_wrap_html = soup.find('div', class_='fact__temp-wrap')
        current_weather = fact_temp_wrap_html.find('a', class_='link fact__basic fact__basic_size_wide day-anchor i-bem').get('aria-label')
        
        current_facts_html = soup.find('div', class_='fact__props')
        for facts in current_facts_html.findAll('span', class_='a11y-hidden'):
            message += facts.text + '\n'

        #today weather:
        today_weather_html = soup.find('div', class_='swiper-container fact__hourly-swiper i-bem fact__hourly-swiper_js_inited swiper-container-horizontal')

        for hour_weather in today_weather_html.find_all('span', class_='a11y-hidden'):
            message += hour_weather.text + '\n'


        #next 10 days weather:
        message += 'Погода на 10 дней вкратце: \n'
        next_10_days_html = soup.find('div', class_='forecast-briefly__days swiper-container swiper-container-horizontal').find('ul')
        
        for day_weather in next_10_days_html.findAll('li'):
            message += (day_weather.find('a', class_='link link_theme_normal text forecast-briefly__day-link i-bem').get('aria-label'))
            message += '\n'
        
        return (True, message)
    
    except Exception as e:
        return (False, f"Произошла ошибка: {e}")

    

def check_valid(soup):
    title = soup.find('title')
    return title.text != 'Ой!'

