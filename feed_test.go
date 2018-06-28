package starling

var feedTC = []struct {
	name string
	mock string
}{
	{
		name: "multiple transactions",
		mock: `
		"feedItems": [
			 {
				  "feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
				  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
				  "amount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "sourceAmount": {
						"currency": "GBP",
						"minorUnits": 32
				  },
				  "direction": "OUT",
				  "transactionTime": "2018-06-28T07:16:28.364Z",
				  "source": "MASTER_CARD",
				  "sourceSubType": "CHIP_AND_PIN",
				  "status": "SETTLED",
				  "counterPartyType": "MERCHANT",
				  "counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
				  "reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
				  "country": "GB",
				  "spendingCategory": "HOLIDAYS"
			 }
			 ]
		}`,
	},
	{
		name: "multiple transactions",
		mock: `{
			"feedItems": [
				 {
					  "feedItemUid": "dbb59f1c-39e6-4558-87ba-11c142965393",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 32
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 32
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.364Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "e6dbe57e-7c23-4015-97a4-4afbbf7faa23",
					  "reference": "ATM 111072\\35 REGENT ST), LONDON\\LONDON\\SW1Y 4ND  00 GBR",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 },
				 {
					  "feedItemUid": "199c2bba-9f4d-4b20-b5df-4de440411b03",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 7
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 7
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.361Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
					  "reference": "MASTERCARD EUROPE      WATERLOO      BEL",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 },
				 {
					  "feedItemUid": "32f8ffc4-d12c-43fe-9d1b-61faf7243143",
					  "categoryUid": "c423ab8d-9a6a-44b2-8db6-ac6000fe58e0",
					  "amount": {
							"currency": "GBP",
							"minorUnits": 24
					  },
					  "sourceAmount": {
							"currency": "GBP",
							"minorUnits": 24
					  },
					  "direction": "OUT",
					  "transactionTime": "2018-06-28T07:16:28.359Z",
					  "source": "MASTER_CARD",
					  "sourceSubType": "CHIP_AND_PIN",
					  "status": "SETTLED",
					  "counterPartyType": "MERCHANT",
					  "counterPartyUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
					  "reference": "MASTERCARD UK MANA\\19TH FLOOR\\LONDON E14\\E14 5NP      GBR",
					  "country": "GB",
					  "spendingCategory": "HOLIDAYS"
				 }
			]
	  }`,
	},
}
