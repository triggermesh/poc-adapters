FROM python:buster

RUN apt-get update

COPY requirements.txt /
RUN pip install -r requirements.txt

COPY copybooktransformation.py /
EXPOSE 8080
CMD ["python", "-u", "/copybooktransformation.py"]
