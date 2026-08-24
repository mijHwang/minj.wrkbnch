import {useState} from "react";
import {shortenUrl, API_BASE} from "../services/urlApi";

export default function ShortenURLDemo(){

    const [longURL, setLongURL] = useState("");
    const [history, setHistory] = useState([]);
    const [loading, setLoading] = useState(false);

    async function handleShorten(e) {
        e.preventDefault();
        setLoading(true);
        const result = await shortenUrl(longURL);
        setHistory((prev) => [...prev, result]);
        setLongURL("");
        setLoading(false);
    }   

    return (
        <div className="demo">

            <h1 className="demo-title">URL shortner demo</h1>
            <p className="demo-intro">
                Paste a long URL below and it's stored in Postgres, with the row's <strong>id</strong> encoded
                into a short, unique code — no randomness, no collisions possible. Click the result to follow
                the real redirect.
            </p>

            <div className="demo-grid">
                
                <div className="panel control-panel">
                    <form className="field-group" onSubmit={handleShorten}>
                        <label>Input URL</label>
                        <input value={longURL} onChange={(e) => setLongURL(e.target.value)} placeholder="URL" required/>
                        <button type="submit" disabled={loading}>
                            {loading ? "Shortening..." : "Shorten"}
                        </button>
                    </form>
                </div>

                <div className="panel history-panel">
                    {history.length === 0 ? (
                        <p className="history-empty">No history of shortening. Shoten a URL to se it appear</p>
                    ): (
                        <div className="history-list">
                        {history.map((record, i) => (
                            <div key={i} className="history-item">
                                <span className="long-url">{record.long_url}</span>    
                                <a href={`${API_BASE}/${record.short_url}`} target="_blank" rel="noreferrer" className="short-url">
                                    {API_BASE}/{record.short_url}
                                </a>
                            </div>
                        ))}
                        </div>
                    )}
                </div>

            </div>
        </div>
    );
}