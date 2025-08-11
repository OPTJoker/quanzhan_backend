"""
AI Agent 开发示例
这是第三个学习模块 - 构建智能代理和工具调用系统
"""

import os
import json
import requests
from datetime import datetime
from typing import List, Dict, Any, Optional
from langchain.agents import initialize_agent, AgentType, Tool
from langchain.llms import OpenAI
from langchain.memory import ConversationBufferMemory
from langchain.tools import BaseTool
from langchain.callbacks.manager import CallbackManagerForToolRun
import mysql.connector
import redis
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

class WeatherTool(BaseTool):
    """天气查询工具"""
    name = "weather_search"
    description = "查询指定城市的当前天气信息。输入应该是城市名称。"
    
    def _run(
        self, 
        query: str, 
        run_manager: Optional[CallbackManagerForToolRun] = None
    ) -> str:
        """执行天气查询"""
        try:
            # 这里应该调用真实的天气API
            # 为了演示，我们返回模拟数据
            return f"当前{query}的天气：晴朗，温度25°C，湿度60%"
        except Exception as e:
            return f"查询天气时出错：{str(e)}"

class DatabaseTool(BaseTool):
    """数据库查询工具"""
    name = "database_query"
    description = "查询数据库中的信息。输入应该是查询的内容描述。"
    
    def __init__(self):
        super().__init__()
        self.conn = mysql.connector.connect(
            host="localhost",
            user="ai_user",
            password="ai_pass",
            database="ai_development"
        )
    
    def _run(
        self, 
        query: str, 
        run_manager: Optional[CallbackManagerForToolRun] = None
    ) -> str:
        """执行数据库查询"""
        try:
            cursor = self.conn.cursor()
            
            # 简化的查询逻辑，实际应用中需要更复杂的SQL生成
            if "聊天" in query or "对话" in query:
                cursor.execute("SELECT COUNT(*) FROM chats")
                count = cursor.fetchone()[0]
                result = f"数据库中共有 {count} 个聊天会话"
            elif "消息" in query:
                cursor.execute("SELECT COUNT(*) FROM messages")
                count = cursor.fetchone()[0]
                result = f"数据库中共有 {count} 条消息"
            else:
                result = "请提供更具体的查询内容"
            
            cursor.close()
            return result
            
        except Exception as e:
            return f"数据库查询出错：{str(e)}"

class CalculatorTool(BaseTool):
    """计算器工具"""
    name = "calculator"
    description = "执行数学计算。输入应该是数学表达式，如 '2+3*4'。"
    
    def _run(
        self, 
        query: str, 
        run_manager: Optional[CallbackManagerForToolRun] = None
    ) -> str:
        """执行数学计算"""
        try:
            # 安全地执行数学表达式
            allowed_chars = set('0123456789+-*/(). ')
            if not all(c in allowed_chars for c in query):
                return "表达式包含不允许的字符"
            
            result = eval(query)
            return f"计算结果：{result}"
        except Exception as e:
            return f"计算出错：{str(e)}"

class AIAgent:
    """AI智能代理"""
    
    def __init__(self):
        self.llm = OpenAI(
            openai_api_key=os.getenv("OPENAI_API_KEY"),
            temperature=0.7
        )
        
        # 初始化工具
        self.tools = [
            WeatherTool(),
            DatabaseTool(),
            CalculatorTool(),
        ]
        
        # 创建记忆
        self.memory = ConversationBufferMemory(
            memory_key="chat_history",
            return_messages=True
        )
        
        # 初始化代理
        self.agent = initialize_agent(
            tools=self.tools,
            llm=self.llm,
            agent=AgentType.CONVERSATIONAL_REACT_DESCRIPTION,
            memory=self.memory,
            verbose=True,
            max_iterations=3,
            early_stopping_method="generate"
        )
        
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
    
    def run(self, user_input: str, session_id: str) -> Dict[str, Any]:
        """运行代理"""
        try:
            # 检查缓存
            cache_key = f"agent_session:{session_id}:query:{hash(user_input)}"
            cached_result = self.redis_conn.get(cache_key)
            
            if cached_result:
                return json.loads(cached_result)
            
            # 执行代理
            response = self.agent.run(user_input)
            
            result = {
                "response": response,
                "session_id": session_id,
                "timestamp": datetime.now().isoformat(),
                "tools_used": self._get_tools_used()
            }
            
            # 缓存结果
            self.redis_conn.setex(
                cache_key,
                1800,  # 30分钟过期
                json.dumps(result, ensure_ascii=False)
            )
            
            # 保存到数据库
            self._save_interaction(session_id, user_input, response)
            
            return result
            
        except Exception as e:
            return {
                "error": str(e),
                "session_id": session_id,
                "timestamp": datetime.now().isoformat()
            }
    
    def _get_tools_used(self) -> List[str]:
        """获取使用的工具列表"""
        # 这里应该从代理的执行历史中提取工具使用信息
        # 简化实现
        return [tool.name for tool in self.tools]
    
    def _save_interaction(self, session_id: str, user_input: str, response: str):
        """保存交互历史"""
        cursor = self.mysql_conn.cursor()
        
        query = """
        INSERT INTO agent_interactions (session_id, user_input, agent_response, created_at) 
        VALUES (%s, %s, %s, NOW())
        """
        
        cursor.execute(query, (session_id, user_input, response))
        self.mysql_conn.commit()
        cursor.close()
    
    def get_session_history(self, session_id: str) -> List[Dict[str, Any]]:
        """获取会话历史"""
        cursor = self.mysql_conn.cursor(dictionary=True)
        
        query = """
        SELECT user_input, agent_response, created_at 
        FROM agent_interactions 
        WHERE session_id = %s 
        ORDER BY created_at ASC
        """
        
        cursor.execute(query, (session_id,))
        results = cursor.fetchall()
        cursor.close()
        
        return results
    
    def add_custom_tool(self, tool: BaseTool):
        """添加自定义工具"""
        self.tools.append(tool)
        
        # 重新初始化代理
        self.agent = initialize_agent(
            tools=self.tools,
            llm=self.llm,
            agent=AgentType.CONVERSATIONAL_REACT_DESCRIPTION,
            memory=self.memory,
            verbose=True,
            max_iterations=3,
            early_stopping_method="generate"
        )

class TaskPlanningAgent(AIAgent):
    """任务规划代理"""
    
    def __init__(self):
        super().__init__()
        
        # 添加任务规划特定的工具
        self.add_custom_tool(TaskManagerTool())
    
    def create_task_plan(self, goal: str) -> Dict[str, Any]:
        """创建任务计划"""
        prompt = f"""
        作为一个任务规划专家，请为以下目标创建详细的任务计划：
        
        目标：{goal}
        
        请提供：
        1. 分解的子任务列表
        2. 每个任务的优先级
        3. 预估的完成时间
        4. 所需的资源和工具
        5. 任务之间的依赖关系
        
        请以结构化的方式回答。
        """
        
        return self.run(prompt, f"task_planning_{hash(goal)}")

class TaskManagerTool(BaseTool):
    """任务管理工具"""
    name = "task_manager"
    description = "管理任务列表，可以创建、更新、删除和查询任务。"
    
    def _run(
        self, 
        query: str, 
        run_manager: Optional[CallbackManagerForToolRun] = None
    ) -> str:
        """执行任务管理操作"""
        # 这里实现任务管理逻辑
        return f"任务管理操作：{query}"

def main():
    """示例用法"""
    print("🤖 AI Agent 开发示例")
    print("=" * 50)
    
    # 创建AI代理
    agent = AIAgent()
    session_id = "demo_session_001"
    
    # 示例交互
    test_queries = [
        "你好，你能做什么？",
        "计算 15 * 23 + 7",
        "查询一下数据库中有多少聊天记录",
        "帮我查一下北京的天气",
        "我想了解你有哪些工具可以使用"
    ]
    
    for query in test_queries:
        print(f"\n👤 用户: {query}")
        result = agent.run(query, session_id)
        
        if "error" in result:
            print(f"❌ 错误: {result['error']}")
        else:
            print(f"🤖 Agent: {result['response']}")
            if result.get('tools_used'):
                print(f"🔧 使用的工具: {', '.join(result['tools_used'])}")
    
    # 任务规划示例
    print("\n" + "=" * 50)
    print("📋 任务规划示例")
    
    planning_agent = TaskPlanningAgent()
    goal = "开发一个在线商城网站"
    
    print(f"\n🎯 目标: {goal}")
    plan_result = planning_agent.create_task_plan(goal)
    
    if "error" not in plan_result:
        print(f"📋 计划: {plan_result['response']}")
    else:
        print(f"❌ 规划失败: {plan_result['error']}")

if __name__ == "__main__":
    main()
