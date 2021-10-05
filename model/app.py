import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
 
app = FastAPI()

class Query(BaseModel):
    body: str  

@app.post("/")
async def process_query(query: Query):
    p = "I've received your query: " + query.body
    return {"body" : p}


if __name__ == '__main__':
    uvicorn.run(app, host='127.0.0.1', port=8000)


