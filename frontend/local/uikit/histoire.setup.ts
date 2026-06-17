import "./historie.css";
import "./variables.css";
import "./variables-dark.css";
import "./fonts.css";
import "./.histoire/assets/style/reset.css";
import "./.histoire/assets/style/doc.css";
import "./.histoire/assets/style/recolor.css"
import { defineSetupVue3 } from "@histoire/plugin-vue";
import TNLifeUIKit from "./index";


export const setupVue3 = defineSetupVue3(({ app }) => {
  app.use(TNLifeUIKit);
})
