import {useState, useEffect} from "react";
import ProjectCard from "../components/ProjectCard";
import CacheThumbnail from "../components/thumbnails/CacheThumbnail";

export default function Hub() {

    //memory allocation for theme, setTheme
    const [theme, setTheme] = useState("light");


    // every time theme changes... the document.documentElement.dataset.theme to theme.
    useEffect(()=>{
        document.documentElement.dataset.theme = theme;
    }, [theme]);
    //bridge between REact's state and the plain CSS

    return (
        <div className = "hub">
            <header className="hub-header">
                <span className="hub-logo">Project Repository</span>
                <div className="hub-actions">
                    <button
                        className="icon-btn"
                        onClick={() => setTheme(theme === "light" ? "dark" : "light")}
                        artial-label="Toggle theme"
                        >
                        {theme === "light"? "☾" : "☀"}
                    </button>

                   <a className="icon-btn"
                    href="https://github.com/mijHwang"
                    target="_blank"
                    rel="noreferrer"
                    aria-label="GitHub"
                    >
                    ⌂
                    </a>
                </div>
            </header>
            <main className="hub-grid">
                <ProjectCard title="LRU Cache Server" thumbnail={<CacheThumbnail />} to="/cache" />
            </main>

            <footer className="hub-footer">
                <p>Copyright © 2026 Min Jun Hwang</p>
            </footer>
        </div>
    );

}

