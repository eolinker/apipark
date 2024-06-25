// denied - 禁用； granted - 拥有权限
// 条件 anyOf/oneOf/anyOf/not
// 维度 backend - 后端的权限字段;

export const PERMISSION_DEFINITION = [
    {
      "system.member.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.member.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.member.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.member.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.member.block": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.department.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.department.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.member.department.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_manager"] }]
        }
      },
      "system.user.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.user.group.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.user.group.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.user.group.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.user.member.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.user.member.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.user_group"] }]
        }
      },
      "system.team.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.team_manager"] }]
        }
      },
      "system.team.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.team_manager"] }]
        }
      },
      "system.team.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.team_manager"] }]
        }
      },
      "system.team.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.team_manager"] }]
        }
      },
      "system.team.self.running": {
        "granted": {
          "anyOf": [{ "backend": ["system.team_manager"] }]
        }
      },
      "system.organization.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.organization_manager"] }]
        }
      },
      "system.organization.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.organization_manager"] }]
        }
      },
      "system.organization.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.organization_manager"] }]
        }
      },
      "system.organization.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.organization_manager"] }]
        }
      },
      "system.role.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.role_manager"] }]
        }
      },
      "system.role.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.role_manager"] }]
        }
      },
      "system.role.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.role_manager"] }]
        }
      },
      "system.role.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.role_manager"] }]
        }
      },
      "system.access.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.system_permission_setting"] }]
        }
      },
      "system.access.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.system_permission_setting"] }]
        }
      },
      "system.access.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.system_permission_setting"] }]
        }
      },
      "system.access.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.system_permission_setting"] }]
        }
      },
      "system.partition.cluster.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cluster.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cluster.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cluster.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cert.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cert.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cert.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.cert.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.partition.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.openapi.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.logRetrieval.self.view":{
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.openapi.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.openapi.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.openapi.self.updateToken": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.openapi.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.auditLog.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["system.environs_setting"] }]
        }
      },
      "system.serviceHub.category.add": {
        "granted": {
          "anyOf": [{ "backend": ["system.service_categories_setting"] }]
        }
      },
      "system.serviceHub.category.edit": {
        "granted": {
          "anyOf": [{ "backend": ["system.service_categories_setting"] }]
        }
      },
      "system.serviceHub.category.delete": {
        "granted": {
          "anyOf": [{ "backend": ["system.service_categories_setting"] }]
        }
      },
      "team.myTeam.system.view": {
        "granted": {
          "anyOf": [{ "backend": ["team.project_manager"] }]
        }
      },
      "team.myTeam.system.add": {
        "granted": {
          "anyOf": [{ "backend": ["team.project_manager"] }]
        }
      },
      "team.mySystem.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["team.project_manager","team.project_view"] }]
        }
      },
      "team.mySystem.self.add": {
        "granted": {
          "anyOf": [{ "backend": ["team.project_manager"] }]
        }
      },
      "team.mySystem.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["team.project_manager"] }]
        }
      },
      "team.myTeam.access.view": {
        "granted": {
          "anyOf": [{ "backend": ["team.team_permission_setting"] }]
        }
      },
      "team.myTeam.access.edit": {
        "granted": {
          "anyOf": [{ "backend": ["team.team_permission_setting"] }]
        }
      },
      "team.myTeam.access.delete": {
        "granted": {
          "anyOf": [{ "backend": ["team.team_permission_setting"] }]
        }
      },
      "team.myTeam.self.view": {
        "granted": {
          "anyOf": [{ "backend": ["team.team_setting"] }]
        }
      },
      "team.myTeam.self.edit": {
        "granted": {
          "anyOf": [{ "backend": ["team.team_setting"] }]
        }
      },
      "team.myTeam.member.view": {
        "granted": {
          "anyOf": [{ "backend": ["team.member_setting"] }]
        }
      },
      "team.myTeam.member.add": {
        "granted": {
          "anyOf": [{ "backend": ["team.member_setting"] }]
        }
      },
      "team.myTeam.member.edit": {
        "granted": {
          "anyOf": [{ "backend": ["team.member_setting"] }]
        }
      },
      "project.mySystem.self.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.project_setting"] }]
        }
      },
      "project.mySystem.member.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.member_setting"] }]
        }
      },
      "project.mySystem.member.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.member_setting"] }]
        }
      },
      "project.mySystem.member.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.member_setting"] }]
        }
      },
      "project.mySystem.api.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager","project.api_view"] }]
        }
      },
      "project.mySystem.api.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager"] }]
        }
      },
      "project.mySystem.api.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager"] }]
        }
      },
      "project.mySystem.api.copy": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager"] }]
        }
      },
      "project.mySystem.api.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager"] }]
        }
      },
      "project.mySystem.api.import": {
        "granted": {
          "anyOf": [{ "backend": ["project.api_manager"] }]
        }
      },
      "project.mySystem.upstream.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.upstream_manager","project.upstream_view"] }]
        }
      },
      "project.mySystem.upstream.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.upstream_manager"] }]
        }
      },
      "project.mySystem.upstream.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.upstream_manager"] }]
        }
      },
      "project.mySystem.upstream.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.upstream_manager"] }]
        }
      },
      "project.mySystem.service.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.service_manager","project.service_view"] }]
        }
      },
      "project.mySystem.service.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.service_manager","project.service_view"] }]
        }
      },
      "project.mySystem.service.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.service_manager"] }]
        }
      },
      "project.mySystem.service.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.service_manager"] }]
        }
      },
      "project.mySystem.service.running": {
        "granted": {
          "anyOf": [{ "backend": ["project.service_manager"] }]
        }
      },
      "project.mySystem.subservice.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_view","project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subservice.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subservice.subscribe": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subservice.viewApproval": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subservice.cancelSubscribe": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subservice.cancelApply": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_apply"] }]
        }
      },
      "project.mySystem.subscriber.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribers_manager"] }]
        }
      },
      "project.mySystem.subscriber.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribers_manager"] }]
        }
      },
      "project.mySystem.subscriber.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribers_manager"] }]
        }
      },
      "project.mySystem.subscribeApproval.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_approval"] }]
        }
      },
      "project.mySystem.subscribeApproval.approval": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_approval"] }]
        }
      },
      "project.mySystem.statistics.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_approval"] }]
        }
      },
      "project.mySystem.topology.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.subscribe_approval"] }]
        }
      },
      "project.mySystem.auth.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.authentication_view","project.authentication_manager"] }]
        }
      },
      "project.mySystem.auth.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.authentication_manager"] }]
        }
      },
      "project.mySystem.auth.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.authentication_manager"] }]
        }
      },
      "project.mySystem.auth.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.authentication_manager"] }]
        }
      },
      "project.mySystem.publish.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.add": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.online": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.stop": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.cancel": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.rollback": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_manager"] }]
        }
      },
      "project.mySystem.publish.approval": {
        "granted": {
          "anyOf": [{ "backend": ["project.publish_approve"] }]
        }
      },
      "project.mySystem.access.view": {
        "granted": {
          "anyOf": [{ "backend": ["project.permission_manager"] }]
        }
      },
      "project.mySystem.access.edit": {
        "granted": {
          "anyOf": [{ "backend": ["project.permission_manager"] }]
        }
      },
      "project.mySystem.access.delete": {
        "granted": {
          "anyOf": [{ "backend": ["project.permission_manager"] }]
        }
      }
    }
  ];