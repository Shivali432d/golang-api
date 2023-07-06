const express = require("express");
const app = express();

const PORT = process.env.PORT || 5000;

const products_routes = require("./routes/products");

app.get("/", (req, res) => {
  res.send("<b>Welcome Shivi!!!!!</b>");
});

// middleware or set router
app.use("/products", products_routes);
app.use("/testing", products_routes);

const start = async () => {
  try {
    app.listen(PORT, () => {
      console.log(`Hola! We're connected on port ${PORT}!`);
    });
  } catch (error) {
    console.log(error);
  }
};

start();
