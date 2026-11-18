use zip::ZipArchive;
use std::fs;
use std::fs::File;
use serde::Serialize;
use serde_json::to_string_pretty;   


#[derive(Serialize)]
struct Report {
    file: String,
    size: u64,
    entry_count: usize,
    entries: Vec<String>,
}

fn analyse_zip(path: &str) -> zip::result::ZipResult<Report> {
    let metadata = fs::metadata(path)?;
    let file = File::open(path)?;
    let mut archive = ZipArchive::new(file)?;
    let mut entries: Vec<String> = Vec::new();

    for i in 0..archive.len() {
        let entry = archive.by_index(i)?;
        let name = entry.name()?.to_string();
        entries.push(name)
    }
    Ok(Report {
        file: path.to_string(),
        size: metadata.len(),
        entry_count: entries.len(),
        entries,
    })

}

fn main() {

    let path = "/tmp/test.zip";
    let report = analyse_zip(path).expect("analysis failed");
    let json = to_string_pretty(&report).expect("serialize failed");
    fs::write("/tmp/report.json", json).expect("cannot write report");
}
