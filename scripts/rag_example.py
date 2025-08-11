"""
RAG (检索增强生成) 实现示例
这是第二个学习模块 - 构建知识库问答系统
"""

import os
import json
from typing import List, Dict, Any
from langchain.vectorstores import Chroma
from langchain.embeddings import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.chains import RetrievalQA
from langchain.llms import OpenAI
from langchain.document_loaders import TextLoader, PyPDFLoader
import mysql.connector
import redis
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

class RAGSystem:
    """RAG知识库问答系统"""
    
    def __init__(self):
        self.embeddings = OpenAIEmbeddings(
            openai_api_key=os.getenv("OPENAI_API_KEY")
        )
        self.llm = OpenAI(
            openai_api_key=os.getenv("OPENAI_API_KEY"),
            temperature=0.7
        )
        self.vectorstore = None
        self.qa_chain = None
        
        # 数据库连接
        self.mysql_conn = mysql.connector.connect(
            host="localhost",
            user="ai_user",
            password="ai_pass",
            database="ai_development"
        )
        
        self.redis_conn = redis.Redis(
            host="localhost",
            port=6379,
            db=0,
            decode_responses=True
        )
    
    def load_documents(self, file_paths: List[str]) -> List[Any]:
        """加载文档"""
        documents = []
        
        for file_path in file_paths:
            if file_path.endswith('.pdf'):
                loader = PyPDFLoader(file_path)
            else:
                loader = TextLoader(file_path, encoding='utf-8')
            
            docs = loader.load()
            documents.extend(docs)
        
        return documents
    
    def split_documents(self, documents: List[Any]) -> List[Any]:
        """分割文档"""
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
        )
        
        return text_splitter.split_documents(documents)
    
    def create_vectorstore(self, documents: List[Any]):
        """创建向量数据库"""
        chunks = self.split_documents(documents)
        
        self.vectorstore = Chroma.from_documents(
            documents=chunks,
            embedding=self.embeddings,
            persist_directory="./data/chroma_db"
        )
        
        # 持久化向量数据库
        self.vectorstore.persist()
        
        # 创建QA链
        self.qa_chain = RetrievalQA.from_chain_type(
            llm=self.llm,
            chain_type="stuff",
            retriever=self.vectorstore.as_retriever(search_kwargs={"k": 3}),
            return_source_documents=True
        )
    
    def load_vectorstore(self):
        """加载已存在的向量数据库"""
        self.vectorstore = Chroma(
            persist_directory="./data/chroma_db",
            embedding_function=self.embeddings
        )
        
        self.qa_chain = RetrievalQA.from_chain_type(
            llm=self.llm,
            chain_type="stuff",
            retriever=self.vectorstore.as_retriever(search_kwargs={"k": 3}),
            return_source_documents=True
        )
    
    def query(self, question: str) -> Dict[str, Any]:
        """查询知识库"""
        if not self.qa_chain:
            return {"error": "RAG系统未初始化"}
        
        # 检查Redis缓存
        cache_key = f"rag_query:{hash(question)}"
        cached_result = self.redis_conn.get(cache_key)
        
        if cached_result:
            return json.loads(cached_result)
        
        # 执行查询
        try:
            result = self.qa_chain({"query": question})
            
            response = {
                "answer": result["result"],
                "sources": [doc.page_content[:200] + "..." for doc in result["source_documents"]],
                "metadata": [doc.metadata for doc in result["source_documents"]]
            }
            
            # 缓存结果
            self.redis_conn.setex(
                cache_key, 
                3600,  # 1小时过期
                json.dumps(response, ensure_ascii=False)
            )
            
            # 保存查询历史到MySQL
            self.save_query_history(question, response["answer"])
            
            return response
            
        except Exception as e:
            return {"error": str(e)}
    
    def save_query_history(self, question: str, answer: str):
        """保存查询历史到数据库"""
        cursor = self.mysql_conn.cursor()
        
        query = """
        INSERT INTO rag_queries (question, answer, created_at) 
        VALUES (%s, %s, NOW())
        """
        
        cursor.execute(query, (question, answer))
        self.mysql_conn.commit()
        cursor.close()
    
    def add_document(self, content: str, title: str, source: str):
        """添加新文档到知识库"""
        # 保存到MySQL
        cursor = self.mysql_conn.cursor()
        
        query = """
        INSERT INTO documents (title, content, source, created_at) 
        VALUES (%s, %s, %s, NOW())
        """
        
        cursor.execute(query, (title, content, source))
        self.mysql_conn.commit()
        cursor.close()
        
        # 更新向量数据库
        if self.vectorstore:
            from langchain.schema import Document
            
            doc = Document(page_content=content, metadata={"title": title, "source": source})
            self.vectorstore.add_documents([doc])
            self.vectorstore.persist()

def main():
    """示例用法"""
    rag_system = RAGSystem()
    
    # 示例1：加载文档并创建知识库
    print("🔄 正在创建知识库...")
    
    # 创建示例文档
    sample_docs = [
        "AI人工智能是计算机科学的一个分支，致力于创建能够执行通常需要人类智能的任务的系统。",
        "机器学习是人工智能的一个子集，它使计算机能够在没有明确编程的情况下学习和改进。",
        "深度学习是机器学习的一个分支，使用人工神经网络来模拟人脑的学习过程。",
        "自然语言处理(NLP)是人工智能的一个领域，专注于计算机与人类语言之间的交互。"
    ]
    
    # 创建临时文档文件
    os.makedirs("./data/docs", exist_ok=True)
    
    for i, content in enumerate(sample_docs):
        with open(f"./data/docs/doc_{i}.txt", "w", encoding="utf-8") as f:
            f.write(content)
    
    # 加载文档
    doc_files = [f"./data/docs/doc_{i}.txt" for i in range(len(sample_docs))]
    documents = rag_system.load_documents(doc_files)
    
    # 创建向量数据库
    rag_system.create_vectorstore(documents)
    
    print("✅ 知识库创建完成!")
    
    # 示例2：查询知识库
    questions = [
        "什么是人工智能？",
        "机器学习和深度学习有什么区别？",
        "NLP是什么？"
    ]
    
    for question in questions:
        print(f"\n❓ 问题: {question}")
        result = rag_system.query(question)
        
        if "error" in result:
            print(f"❌ 错误: {result['error']}")
        else:
            print(f"🤖 回答: {result['answer']}")
            print(f"📚 来源: {len(result['sources'])} 个相关文档片段")

if __name__ == "__main__":
    main()
