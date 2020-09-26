1. chatserver在项目中的目录结构
ChatServer
   |--------src
             |------bll           // 业务逻辑层
                   |----handle    // 消息注册及处理
                   |----roomMgr   // 房间管理
                   |----sensitive // 敏感字
             |------model         // 模型定义
             |------ws            // socket网络处理
