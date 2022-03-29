FROM python:buster

RUN apt-get update

COPY requirements.txt /
RUN pip install -r requirements.txt

COPY csnfeventtransformation.py /
COPY onug_decorator.py /

EXPOSE 8080
CMD ["python", "-u", "/csnfeventtransformation.py"]
