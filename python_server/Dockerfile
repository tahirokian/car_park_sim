FROM python:3.10-alpine

COPY ./requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt

RUN mkdir /app
WORKDIR /app

COPY ./app /app

EXPOSE 40000

CMD ["python", "app.py"]
