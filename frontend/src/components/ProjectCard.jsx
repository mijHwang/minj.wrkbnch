import { Link } from "react-router-dom";

//Link is a react route version of <a>
//navigate without full page reload


export default function ProjectCard({title, thumbnail, to, href}){

    const content = (
        <>
            <div className="card-thumb">{thumbnail}</div>
            <div className="card-label">{title}</div>
        </>
    );


    return to ? (
        <Link to={to} className="project-card">{content}</Link>
    ) : (
        <a href={href} className="project-card">{content}</a>
    );

}