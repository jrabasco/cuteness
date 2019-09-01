FROM scratch
COPY cuteness /app/cuteness
COPY cuteness.html /app/cuteness.html
COPY authenticate.html /app/authenticate.html
WORKDIR /app
CMD ["./cuteness", "-images", "/assets/images/cuteness", "-texts", "/texts"]
