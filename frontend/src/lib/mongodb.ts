import { MongoClient } from "mongodb";

const uri = process.env.MONGODB_URI;
if (!uri) {
  throw new Error("Please define the MONGODB_URI environment variable");
}

let clientPromise: Promise<MongoClient>;

const client = new MongoClient(uri);
clientPromise = client.connect();

export default clientPromise;
