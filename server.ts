import express from "express"
import {RateLimiter} from "./utils.js"

const app = express()
const PORT = 3000

app.get("/arcane", async(req, res) => {
  console.log("arcane")
  await RateLimiter(req, res)
})

app.listen(PORT, () => {
  console.log(`LISTENING AT PORT `, PORT)
})
 

interface Commit {
  Author: string
  Tree: Hash
  Hash: Hash
}

interface Hash {
  Hash: string
}
