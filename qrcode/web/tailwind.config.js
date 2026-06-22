/** @type {import('tailwindcss').Config} */

export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,vue}"],
  theme: {
    container: {
      center: true,
    },
    extend: {
      colors: {
        base: {
          DEFAULT: "#0E0E11",
          panel: "#16161B",
          elevated: "#1C1C24",
          hover: "#22222B",
          border: "#26262E",
          "border-strong": "#33333D",
        },
        ink: {
          DEFAULT: "#F4F4F5",
          secondary: "#A1A1AA",
          muted: "#71717A",
          faint: "#52525B",
        },
        lime: {
          DEFAULT: "#D4FF3A",
          soft: "#B8E02E",
          dim: "#7A9420",
          glow: "rgba(212,255,58,0.35)",
        },
        amber: {
          DEFAULT: "#FFB454",
          dim: "#8A5C2A",
          glow: "rgba(255,180,84,0.35)",
        },
      },
      fontFamily: {
        display: ['"Syne"', "sans-serif"],
        sans: ['"IBM Plex Sans"', "system-ui", "sans-serif"],
        mono: ['"JetBrains Mono"', "monospace"],
      },
      borderRadius: {
        xl: "12px",
        "2xl": "16px",
      },
      boxShadow: {
        glow: "0 0 0 1px rgba(212,255,58,0.4), 0 0 24px rgba(212,255,58,0.18)",
        "glow-amber": "0 0 0 1px rgba(255,180,84,0.4), 0 0 24px rgba(255,180,84,0.18)",
        panel: "0 1px 0 rgba(255,255,255,0.03) inset, 0 8px 24px rgba(0,0,0,0.35)",
      },
      backgroundImage: {
        "grid-faint":
          "linear-gradient(rgba(255,255,255,0.025) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.025) 1px, transparent 1px)",
        "checker":
          "conic-gradient(#1c1c24 25%, #16161b 0 50%, #1c1c24 0 75%, #16161b 0)",
      },
      backgroundSize: {
        grid: "32px 32px",
        checker: "16px 16px",
      },
      keyframes: {
        "fade-up": {
          "0%": { opacity: "0", transform: "translateY(8px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        scan: {
          "0%": { transform: "translateY(-100%)", opacity: "0" },
          "20%": { opacity: "1" },
          "100%": { transform: "translateY(2400%)", opacity: "0" },
        },
        "pulse-dot": {
          "0%,100%": { opacity: "1" },
          "50%": { opacity: "0.35" },
        },
      },
      animation: {
        "fade-up": "fade-up 0.5s cubic-bezier(0.16,1,0.3,1) both",
        scan: "scan 1.4s ease-in-out",
        "pulse-dot": "pulse-dot 1.6s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};
