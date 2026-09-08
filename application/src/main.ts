import { mount } from "svelte"
import "./main.css"
import Main from "./Main.svelte"

export default mount(Main, {
  target: document.getElementById("app")!,
})

