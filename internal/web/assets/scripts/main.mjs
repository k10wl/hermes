import { initConnectionIndicator } from "./ui/connection-indicator.mjs";
import { Publisher } from "./utils/publisher.mjs";

initConnectionIndicator();

const hello = new Publisher("hello world");
hello.notify();
