import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

const config: Config = {
    title: "Velocity-ci",
    tagline: "Streamline your CI/CD pipeline",
    favicon: "img/favicon.ico",
    url: "https://docs.velocity-ci.com",
    baseUrl: "/",
    onBrokenLinks: "throw",
    onBrokenMarkdownLinks: "warn",
    i18n: {
        defaultLocale: "en",
        locales: ["en"],
    },
    presets: [
        [
            "classic",
            {
                docs: {
                    sidebarPath: "./sidebars.ts",
                    editUrl:
                        "https://github.com/zackarysantana/velocity/tree/main/docs",
                },
                blog: {
                    showReadingTime: true,
                    editUrl:
                        "https://github.com/zackarysantana/velocity/tree/main/docs",
                },
                theme: {
                    customCss: "./src/css/custom.css",
                },
            } satisfies Preset.Options,
        ],
    ],
    plugins: [
        [
            "@docusaurus/plugin-content-blog",
            {
                id: "changelog",
                routeBasePath: "changelog",
                path: "./changelog",
                blogSidebarTitle: "Changelog",
            },
        ],
    ],
    themeConfig: {
        // Replace with your project's social card
        image: "img/docusaurus-social-card.jpg",
        navbar: {
            title: "Velocity-ci",
            logo: {
                alt: "My Site Logo",
                src: "img/android-chrome-512x512.png",
            },
            items: [
                {
                    type: "docSidebar",
                    sidebarId: "tutorialSidebar",
                    position: "left",
                    label: "Docs",
                },
                { to: "/changelog", label: "Changelog", position: "left" },
                { to: "/blog", label: "Blog", position: "left" },
                {
                    href: "https://github.com/zackarysantana/velocity",
                    label: "GitHub",
                    position: "right",
                },
            ],
        },
        footer: {
            style: "dark",
            links: [
                {
                    title: "Other Project",
                    items: [
                        {
                            label: "Autumn Changelog",
                            href: "https://autumn-cl.com",
                        },
                        {
                            label: "How's It",
                            href: "https://howsit.dev",
                        },
                        {
                            label: "Portfolio",
                            href: "https://zackaryjamessantana.com",
                        },
                    ],
                },
                {
                    title: "Community",
                    items: [
                        {
                            label: "GitHub",
                            href: "https://github.com/zackarysantana/velocity",
                        },
                    ],
                },
                {
                    title: "More",
                    items: [
                        {
                            label: "Changelog",
                            to: "/changelog",
                        },
                        {
                            label: "Blog",
                            to: "/blog",
                        },
                    ],
                },
            ],
            copyright: `Copyright © ${new Date().getFullYear()} My Project, Inc. Built with Docusaurus.`,
        },
        prism: {
            theme: prismThemes.github,
            darkTheme: prismThemes.dracula,
        },
    } satisfies Preset.ThemeConfig,
};

export default config;
