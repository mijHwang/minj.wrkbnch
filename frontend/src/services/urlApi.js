const API_BASE = import.meta.env.VITE_SHORTENER_API_URL || "http://localhost:8081"

export async function shortenUrl (longUrl){
    const res = await fetch(`${API_BASE}/shorten`,{
        method: "POST",
        headers: {"Content-Type": "application/json" },
        body: JSON.stringify({url: longUrl}),
    });
    return res.json();
}

/* export async function handleRedirect (shortURL){

    const res = await fetch(`/get?=${encodeURIComponent(shortURL)}`)
    return res.json();
} */