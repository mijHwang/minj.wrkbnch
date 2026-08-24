import { useState, useEffect } from "react";
import { setItem, getItem, getAll } from "../services/cacheApi";
import "../assets/cache-demo.css";

export default function CacheDemo() {
  const [items, setItems] = useState([]);
  const [setKey, setSetKey] = useState("");
  const [setValue, setSetValue] = useState("");
  const [ttl, setTtl] = useState(30);
  const [getKey, setGetKey] = useState("");
  const [getResult, setGetResult] = useState(null);
  const [hits, setHits] = useState(0);
  const [misses, setMisses] = useState(0);
  const [setLoading, setSetLoading] = useState(false);
  const [getLoading, setGetLoading] = useState(false);

  async function refresh() {
    const data = await getAll();
    setItems(data);
  }

  useEffect(() => {
    refresh();
    const interval = setInterval(refresh, 2000);
    return () => clearInterval(interval);
  }, []);

  async function handleSet(e) {
    e.preventDefault();
    setSetLoading(true);
    await setItem(setKey, setValue, Number(ttl));
    setSetKey("");
    setSetValue("");
    await refresh();
    setSetLoading(false);
  }

  async function handleGet(e) {
    e.preventDefault();
    setGetLoading(true);
    const result = await getItem(getKey);
    setGetResult(result);
    if (result.found) setHits((h) => h + 1);
    else setMisses((m) => m + 1);
    setGetLoading(false);
  }

  return (
    <div className="demo">
        <h1 className="demo-title">LRU cache demo</h1>
        <p className="demo-intro">
            An <strong>LRU (least recently used) cache</strong> holds a fixed number of items — 
            when it's full, adding a new one evicts whoever hasn't been touched the longest. 
            <strong> Set</strong> a key/value below, then <strong>Get</strong> it back to see it jump 
            to the front. Fill past 5 items to watch the oldest get evicted.
    </p>

    <div className="demo-grid">
        <div className="panel viz-panel">
          {items.length === 0 ? (
            <p className="viz-empty">Cache is empty — set a key to see it appear</p>
          ) : (
            <div className="node-row">
              {items.map((item, i) => (
                <div key={item.key} style={{ display: "flex", alignItems: "center" }}>
                  <div className={`node-box${i === 0 ? " mru" : ""}`}>
                    <div className="node-key">{item.key}</div>
                    <div className="node-val">{item.value}</div>
                    <div className="node-ttl">{item.ttl_remaining_seconds}s left</div>
                  </div>
                  {i < items.length - 1 && <span className="node-arrow">→</span>}
                </div>
              ))}
            </div>
          )}
          <div className="stats-row">
            <span className="stat-hit">Hits: {hits}</span>
            <span className="stat-miss">Misses: {misses}</span>
          </div>
          <p style={{ fontSize: "12px", opacity: 0.5, marginTop: "24px" }}>
            Note: this cache is shared across everyone currently viewing this demo — 
            it's a single in-memory instance, not session-isolated. Resets on redeploy.   
            <strong> Server might need time to wake up on the first try!</strong>
            <br />
            (Also on the TO-DO list. A reset button.)
            </p>
        </div>

        <div className="panel control-panel">
          <form className="field-group" onSubmit={handleSet}>
            <label>Set a key</label>
            <input value={setKey} onChange={(e) => setSetKey(e.target.value)} placeholder="key" required />
            <input value={setValue} onChange={(e) => setSetValue(e.target.value)} placeholder="value" required />
            <input type="number" value={ttl} onChange={(e) => setTtl(e.target.value)} placeholder="ttl seconds" />
            <button type="submit" disabled ={setLoading}>
              {setLoading ? "Setting the key..." : "Set"}
              </button>
          </form>

          <form className="field-group" onSubmit={handleGet}>
            <label>Get a key</label>
            <input value={getKey} onChange={(e) => setGetKey(e.target.value)} placeholder="key to get" required />
            <button type="submit" disabled = {getLoading}>
              {getLoading ? "Getting the key..." : "Get"}
            </button>
          </form>

          {getResult && (
            <div className={`get-result ${getResult.found ? "found" : "missing"}`}>
              {getResult.found ? `Found: ${getResult.value}` : "Not found"}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}