import { defineConfig } from "vitepress";

export default defineConfig({
  title: "🚀 Buzzlink",
  description: "A cli tool for sharing files easily",
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: "Home", link: "/" },
      { text: "Installation", link: "/installation" },
      { text: "Docs", link: "/docs" },
      {
        text: "Releases",
        link: "https://github.com/scifisatan/buzzlink/releases/",
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/scifisatan/buzzlink" },
    ],
  },
  lastUpdated: true,
  markdown: {
    lineNumbers: true,
  },
});
