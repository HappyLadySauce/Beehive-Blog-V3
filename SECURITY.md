# Security Policy

## Supported Versions

| Version | Supported |
| --- | --- |
| main | yes |
| other branches | no |

## Reporting a Vulnerability

如果你发现潜在安全漏洞，请不要公开提交 Issue。  
请通过以下渠道私下披露：

- 13452552349@163.com

## Response SLA

- 24 小时内：确认已收到报告
- 72 小时内：完成初步风险评估并给出处理计划
- 修复完成后：通知报告人并给出影响范围与缓解建议

## Disclosure Rules

- 在官方修复发布前，请勿公开漏洞细节
- 报告内容建议包含：
  - 影响范围
  - 复现步骤
  - 利用条件
  - 风险等级建议

## Security Best Practices

- 不在日志中输出密钥、令牌、密码等敏感字段
- 对外暴露接口要执行输入校验与访问控制
- 依赖升级优先处理高危漏洞
