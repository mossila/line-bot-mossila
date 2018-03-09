from googlefinance import getQuotes
import json
symbol = 'BBL'
print(json.dumps(getQuotes('SET:'+symbol), indent=2))

