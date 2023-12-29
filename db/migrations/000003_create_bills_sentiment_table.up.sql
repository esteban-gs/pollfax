CREATE TABLE bills_sentiments (
  bill_id INTEGER REFERENCES bills(id),
  sentiment_id INTEGER REFERENCES sentiments(id),
  CONSTRAINT bills_sentiments_pk PRIMARY KEY(bill_id,sentiment_id)
);
