import optparse

parser = optparse.OptionParser()
parser.add_option('-u', "--user", action="store", dest="mail_address", help="gmail address")
parser.add_option('-p', "--password", action="store", dest="password", help="gmail password")

options, args = parser.parse_args()


mail_address = options.mail_address
password =  options.password

from selenium import webdriver
import time
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.support.wait import WebDriverWait

options = Options()
options.headless = False


driver = webdriver.Firefox(options=options)

url = 'https://mester.inf.elte.hu/faces/login.xhtml'
driver.get(url)
url = driver.find_element_by_tag_name('a').get_attribute('href')
driver.get(url)

element = WebDriverWait(driver, 15).until(
    lambda x: x.find_element_by_id("identifierId"))

driver.find_element_by_id("identifierId").send_keys(mail_address)
driver.find_element_by_id("identifierNext").click()

time.sleep(10)

driver.find_element_by_name("password").send_keys(password)
element = driver.find_element_by_id('passwordNext')
driver.execute_script("arguments[0].click();", element)

element = WebDriverWait(driver, 15).until(
    lambda x: x.find_element_by_name("j_idt9:j_idt11"))
driver.find_element_by_name("j_idt9:j_idt11").click()

time.sleep(10)

cookies_list = driver.get_cookies()
cookies_dict = {}
for cookie in cookies_list:
    cookies_dict[cookie['name']] = cookie['value']

print(cookies_dict['JSESSIONID'])
