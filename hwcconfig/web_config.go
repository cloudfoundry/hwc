package hwcconfig

import (
	"os"
	"text/template"
)

func (c *HwcConfig) generateWebConfig() error {
	file, err := os.Create(c.WebConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("webconfig").Parse(webConfigTemplate))
	if err := tmpl.Execute(file, c); err != nil {
		return err
	}
	return nil
}

const webConfigTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!-- the root web configuration file -->
<configuration>
    <!--
        Using a location directive with a missing path attribute
        scopes the configuration to the entire machine.  If used in
        conjunction with allowOverride="false", it can be used to
        prevent configuration from being altered on the machine

        Administrators that want to restrict permissions granted to
        web applications should change the default Trust level and ensure
        that overrides are not allowed
    -->
    <location allowOverride="true">
        <system.web>
            <securityPolicy>
                <trustLevel name="Full" policyFile="internal" />
                <trustLevel name="High" policyFile="web_hightrust.config" />
                <trustLevel name="Medium" policyFile="web_mediumtrust.config" />
                <trustLevel name="Low"  policyFile="web_lowtrust.config" />
                <trustLevel name="Minimal" policyFile="web_minimaltrust.config" />
            </securityPolicy>
            <trust level="Full" originUrl="" />
            <fullTrustAssemblies>
                <add
                    assemblyName="Microsoft.VisualStudio.Enterprise.AspNetHelper"
                    version="11.0.0.0"
                    publicKey="002400000480000094000000060200000024000052534131000400000100010007D1FA57C4AED9F0A32E84AA0FAEFD0DE9E8FD6AEC8F87FB03766C834C99921EB23BE79AD9D5DCC1DD9AD236132102900B723CF980957FC4E177108FC607774F29E8320E92EA05ECE4E821C0A5EFE8F1645C4C0C93C1AB99285D622CAA652C1DFAD63D745D6F2DE5F17E5EAF0FC4963D261C8A12436518206DC093344D5AD293"
                    />
                <add
                    assemblyName="Microsoft.VisualStudio.Web"
                    version="11.0.0.0"
                    publicKey="002400000480000094000000060200000024000052534131000400000100010007D1FA57C4AED9F0A32E84AA0FAEFD0DE9E8FD6AEC8F87FB03766C834C99921EB23BE79AD9D5DCC1DD9AD236132102900B723CF980957FC4E177108FC607774F29E8320E92EA05ECE4E821C0A5EFE8F1645C4C0C93C1AB99285D622CAA652C1DFAD63D745D6F2DE5F17E5EAF0FC4963D261C8A12436518206DC093344D5AD293"
                    />
                <add
                    assemblyName="Microsoft.Web.Infrastructure"
                    version="1.0.0.0"
                    publicKey="0024000004800000940000000602000000240000525341310004000001000100B5FC90E7027F67871E773A8FDE8938C81DD402BA65B9201D60593E96C492651E889CC13F1415EBB53FAC1131AE0BD333C5EE6021672D9718EA31A8AEBD0DA0072F25D87DBA6FC90FFD598ED4DA35E44C398C454307E8E33B8426143DAEC9F596836F97C8F74750E5975C64E2189F45DEF46B2A2B1247ADC3652BF5C308055DA9"
                    />
                <add
                    assemblyName="System.Data.SqlServerCe"
                    version="4.0.0.0"
                    publicKey="0024000004800000940000000602000000240000525341310004000001000100272736ad6e5f9586bac2d531eabc3acc666c2f8ec879fa94f8f7b0327d2ff2ed523448f83c3d5c5dd2dfc7bc99c5286b2c125117bf5cbe242b9d41750732b2bdffe649c6efb8e5526d526fdd130095ecdb7bf210809c6cdad8824faa9ac0310ac3cba2aa0523567b2dfa7fe250b30facbd62d4ec99b94ac47c7d3b28f1f6e4c8"
                    />
            </fullTrustAssemblies>
            <partialTrustVisibleAssemblies />
        </system.web>
    </location>

    <system.codedom>
        <compilers>
            <compiler language="c#;cs;csharp" extension=".cs" warningLevel="4" type="Microsoft.CSharp.CSharpCodeProvider, System, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089">
                <providerOption name="CompilerVersion" value="v4.0"/>
                <providerOption name="WarnAsError" value="false"/>
            </compiler>
            <compiler language="vb;vbs;visualbasic;vbscript" extension=".vb" warningLevel="4" type="Microsoft.VisualBasic.VBCodeProvider, System, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089">
                <providerOption name="CompilerVersion" value="v4.0"/>
                <providerOption name="OptionInfer" value="true"/>
                <providerOption name="WarnAsError" value="false"/>
            </compiler>
        </compilers>
    </system.codedom>

    <system.net>
        <defaultProxy>
            <proxy usesystemdefault="true" />
        </defaultProxy>
    </system.net>
    <system.web>
        <authorization>
            <allow users="*" />
        </authorization>

        <browserCaps userAgentCacheKeyLength="64">
            <result type="System.Web.Mobile.MobileCapabilities, System.Web.Mobile, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
        </browserCaps>

        <clientTarget>
            <add alias="uplevel" userAgent="Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)" />
            <add alias="downlevel" userAgent="Generic Downlevel" />
        </clientTarget>

				<compilation tempDirectory="{{.TempDirectory}}">
            <assemblies>
                <add assembly="mscorlib" />
                <add assembly="Microsoft.CSharp, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Configuration, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System.Web, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System.Data, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Web.Services, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System.Xml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Drawing, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System.EnterpriseServices, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add assembly="System.IdentityModel, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Runtime.Serialization, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.ServiceModel, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.ServiceModel.Web, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.Activities, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.ServiceModel.Activities, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.WorkflowServices, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.Core, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
                <add assembly="System.Data.DataSetExtensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.Xml.Linq, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
                <add assembly="System.ComponentModel.DataAnnotations, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.Web.DynamicData, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add assembly="System.Web.ApplicationServices, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
                <add assembly="*" />
            </assemblies>
            <buildProviders>
                <add extension=".aspx" type="System.Web.Compilation.PageBuildProvider" />
                <add extension=".ascx" type="System.Web.Compilation.UserControlBuildProvider" />
                <add extension=".master" type="System.Web.Compilation.MasterPageBuildProvider" />
                <add extension=".asmx" type="System.Web.Compilation.WebServiceBuildProvider" />
                <add extension=".ashx" type="System.Web.Compilation.WebHandlerBuildProvider" />
                <add extension=".soap" type="System.Web.Compilation.WebServiceBuildProvider" />
                <add extension=".resx" type="System.Web.Compilation.ResXBuildProvider" />
                <add extension=".resources" type="System.Web.Compilation.ResourcesBuildProvider" />
                <add extension=".wsdl" type="System.Web.Compilation.WsdlBuildProvider" />
                <add extension=".xsd" type="System.Web.Compilation.XsdBuildProvider" />
                <add extension=".js" type="System.Web.Compilation.ForceCopyBuildProvider" />
                <add extension=".lic" type="System.Web.Compilation.IgnoreFileBuildProvider" />
                <add extension=".licx" type="System.Web.Compilation.IgnoreFileBuildProvider" />
                <add extension=".exclude" type="System.Web.Compilation.IgnoreFileBuildProvider" />
                <add extension=".refresh" type="System.Web.Compilation.IgnoreFileBuildProvider" />
                <add extension=".edmx" type="System.Data.Entity.Design.AspNet.EntityDesignerBuildProvider" />
                <add extension=".xoml" type="System.ServiceModel.Activation.WorkflowServiceBuildProvider, System.WorkflowServices, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add extension=".svc" type="System.ServiceModel.Activation.ServiceBuildProvider, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
                <add extension=".xamlx" type="System.Xaml.Hosting.XamlBuildProvider, System.Xaml.Hosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
            </buildProviders>
            <expressionBuilders>
                <add expressionPrefix="Resources" type="System.Web.Compilation.ResourceExpressionBuilder" />
                <add expressionPrefix="ConnectionStrings" type="System.Web.Compilation.ConnectionStringsExpressionBuilder" />
                <add expressionPrefix="AppSettings" type="System.Web.Compilation.AppSettingsExpressionBuilder" />
                <add expressionPrefix="RouteUrl" type="System.Web.Compilation.RouteUrlExpressionBuilder"/>
                <add expressionPrefix="RouteValue" type="System.Web.Compilation.RouteValueExpressionBuilder"/>
            </expressionBuilders>
            <folderLevelBuildProviders>
                <add name="DataServiceBuildProvider" type="System.Data.Services.BuildProvider.DataServiceBuildProvider, System.Data.Services.Design, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"/>
            </folderLevelBuildProviders>
        </compilation>

        <healthMonitoring>
            <bufferModes>
                <add name="Critical Notification" maxBufferSize="100" maxFlushSize="20"
                    urgentFlushThreshold="1" regularFlushInterval="Infinite" urgentFlushInterval="00:01:00"
                    maxBufferThreads="1" />
                <add name="Notification" maxBufferSize="300" maxFlushSize="20"
                    urgentFlushThreshold="1" regularFlushInterval="Infinite" urgentFlushInterval="00:01:00"
                    maxBufferThreads="1" />
                <add name="Analysis" maxBufferSize="1000" maxFlushSize="100"
                    urgentFlushThreshold="100" regularFlushInterval="00:05:00"
                    urgentFlushInterval="00:01:00" maxBufferThreads="1" />
                <add name="Logging" maxBufferSize="1000" maxFlushSize="200" urgentFlushThreshold="800"
                    regularFlushInterval="00:30:00" urgentFlushInterval="00:05:00"
                    maxBufferThreads="1" />
            </bufferModes>

            <providers>
                <add name="EventLogProvider" type="System.Web.Management.EventLogWebEventProvider,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a" />
                <add connectionStringName="LocalSqlServer" maxEventDetailsLength="1073741823"
                    buffer="false" bufferMode="Notification" name="SqlWebEventProvider"
                    type="System.Web.Management.SqlWebEventProvider,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a" />
                <add name="WmiWebEventProvider" type="System.Web.Management.WmiWebEventProvider,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a" />
            </providers>

            <profiles>
                <add name="Default" minInstances="1" maxLimit="Infinite" minInterval="00:01:00"
                    custom="" />
                <add name="Critical" minInstances="1" maxLimit="Infinite" minInterval="00:00:00"
                    custom="" />
            </profiles>

            <rules>
                <add name="All Errors Default" eventName="All Errors" provider="EventLogProvider"
                    profile="Default" minInstances="1" maxLimit="Infinite" minInterval="00:01:00"
                    custom="" />
                <add name="Failure Audits Default" eventName="Failure Audits"
                    provider="EventLogProvider" profile="Default" minInstances="1"
                    maxLimit="Infinite" minInterval="00:01:00" custom="" />
            </rules>

            <eventMappings>
                <add name="All Events" type="System.Web.Management.WebBaseEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Heartbeats" type="System.Web.Management.WebHeartbeatEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Application Lifetime Events" type="System.Web.Management.WebApplicationLifetimeEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Request Processing Events" type="System.Web.Management.WebRequestEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="All Errors" type="System.Web.Management.WebBaseErrorEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Infrastructure Errors" type="System.Web.Management.WebErrorEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Request Processing Errors" type="System.Web.Management.WebRequestErrorEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="All Audits" type="System.Web.Management.WebAuditEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Failure Audits" type="System.Web.Management.WebFailureAuditEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
                <add name="Success Audits" type="System.Web.Management.WebSuccessAuditEvent,System.Web,Version=4.0.0.0,Culture=neutral,PublicKeyToken=b03f5f7f11d50a3a"
                    startEventCode="0" endEventCode="2147483647" />
            </eventMappings>

        </healthMonitoring>

        <httpHandlers>
            <add path="eurl.axd" verb="*" type="System.Web.HttpNotFoundHandler" validate="True" />
            <add path="trace.axd" verb="*" type="System.Web.Handlers.TraceHandler" validate="True" />
            <add path="WebResource.axd" verb="GET" type="System.Web.Handlers.AssemblyResourceLoader" validate="True" />
            <add verb="*" path="*_AppService.axd" type="System.Web.Script.Services.ScriptHandlerFactory, System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False" />
            <add verb="GET,HEAD" path="ScriptResource.axd" type="System.Web.Handlers.ScriptResourceHandler, System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False"/>
            <add path="*.axd" verb="*" type="System.Web.HttpNotFoundHandler" validate="True" />
            <add path="*.aspx" verb="*" type="System.Web.UI.PageHandlerFactory" validate="True" />
            <add path="*.ashx" verb="*" type="System.Web.UI.SimpleHandlerFactory" validate="True" />
            <add path="*.asmx" verb="*" type="System.Web.Script.Services.ScriptHandlerFactory, System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False" />
            <add path="*.rem" verb="*" type="System.Runtime.Remoting.Channels.Http.HttpRemotingHandlerFactory, System.Runtime.Remoting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" validate="False" />
            <add path="*.soap" verb="*" type="System.Runtime.Remoting.Channels.Http.HttpRemotingHandlerFactory, System.Runtime.Remoting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" validate="False" />
            <add path="*.asax" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.ascx" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.master" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.skin" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.browser" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.sitemap" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.dll.config" verb="GET,HEAD" type="System.Web.StaticFileHandler" validate="True" />
            <add path="*.exe.config" verb="GET,HEAD" type="System.Web.StaticFileHandler" validate="True" />
            <add path="*.config" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.cs" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.csproj" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.vb" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.vbproj" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.webinfo" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.licx" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.resx" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.resources" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.mdb" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.vjsproj" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.java" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.jsl" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.ldb" verb="*" type="System.Web.HttpForbiddenHandler"  validate="True" />
            <add path="*.ad" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.dd" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.ldd" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.sd" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.cd" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.adprototype" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.lddprototype" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.sdm" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.sdmDocument" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.mdf" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.ldf" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.exclude" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.refresh" verb="*" type="System.Web.HttpForbiddenHandler" validate="True" />
            <add path="*.svc" verb="*" type="System.ServiceModel.Activation.HttpHandler, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False"/>
            <add path="*.rules" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*.xoml" verb="*" type="System.ServiceModel.Activation.HttpHandler, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False"/>
            <add path="*.xamlx" verb="*" type="System.Xaml.Hosting.XamlHttpHandlerFactory, System.Xaml.Hosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" validate="False"/>
            <add path="*.aspq" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*.cshtm" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*.cshtml" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*.vbhtm" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*.vbhtml" verb="*" type="System.Web.HttpForbiddenHandler" validate="True"/>
            <add path="*" verb="GET,HEAD,POST" type="System.Web.DefaultHttpHandler" validate="True" />
            <add path="*" verb="*" type="System.Web.HttpMethodNotAllowedHandler" validate="True" />
        </httpHandlers>

        <httpModules>
            <add name="OutputCache" type="System.Web.Caching.OutputCacheModule" />
            <add name="Session" type="System.Web.SessionState.SessionStateModule" />
            <add name="WindowsAuthentication" type="System.Web.Security.WindowsAuthenticationModule" />
            <add name="FormsAuthentication" type="System.Web.Security.FormsAuthenticationModule" />
            <add name="PassportAuthentication" type="System.Web.Security.PassportAuthenticationModule" />
            <add name="RoleManager" type="System.Web.Security.RoleManagerModule" />
            <add name="UrlAuthorization" type="System.Web.Security.UrlAuthorizationModule" />
            <add name="FileAuthorization" type="System.Web.Security.FileAuthorizationModule" />
            <add name="AnonymousIdentification" type="System.Web.Security.AnonymousIdentificationModule" />
            <add name="Profile" type="System.Web.Profile.ProfileModule" />
            <add name="ErrorHandlerModule" type="System.Web.Mobile.ErrorHandlerModule, System.Web.Mobile, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
            <add name="ServiceModel" type="System.ServiceModel.Activation.HttpModule, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
            <add name="UrlRoutingModule-4.0" type="System.Web.Routing.UrlRoutingModule" />
            <add name="ScriptModule-4.0" type="System.Web.Handlers.ScriptModule, System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
        </httpModules>

        <mobileControls sessionStateHistorySize="6" cookielessDataDictionaryType="System.Web.Mobile.CookielessData">
            <device name="XhtmlDeviceAdapters"
                predicateClass="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlPageAdapter"
                predicateMethod="DeviceQualifies"
                pageAdapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlPageAdapter">

                <control name="System.Web.UI.MobileControls.Panel"             adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlPanelAdapter" />
                <control name="System.Web.UI.MobileControls.Form"              adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlFormAdapter" />
                <control name="System.Web.UI.MobileControls.TextBox"           adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlTextBoxAdapter" />
                <control name="System.Web.UI.MobileControls.Label"             adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlLabelAdapter" />
                <control name="System.Web.UI.MobileControls.LiteralText"       adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlLiteralTextAdapter" />
                <control name="System.Web.UI.MobileControls.Link"              adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlLinkAdapter" />
                <control name="System.Web.UI.MobileControls.Command"           adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlCommandAdapter" />
                <control name="System.Web.UI.MobileControls.PhoneCall"         adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlPhoneCallAdapter" />
                <control name="System.Web.UI.MobileControls.List"              adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlListAdapter" />
                <control name="System.Web.UI.MobileControls.SelectionList"     adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlSelectionListAdapter" />
                <control name="System.Web.UI.MobileControls.ObjectList"        adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlObjectListAdapter" />
                <control name="System.Web.UI.MobileControls.Image"             adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlImageAdapter" />
                <control name="System.Web.UI.MobileControls.ValidationSummary" adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlValidationSummaryAdapter" />
                <control name="System.Web.UI.MobileControls.Calendar"          adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlCalendarAdapter" />
                <control name="System.Web.UI.MobileControls.TextView"          adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlTextViewAdapter" />
                <control name="System.Web.UI.MobileControls.MobileControl"     adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlControlAdapter" />
                <control name="System.Web.UI.MobileControls.BaseValidator"     adapter="System.Web.UI.MobileControls.Adapters.XhtmlAdapters.XhtmlValidatorAdapter" />
            </device>
            <device name="HtmlDeviceAdapters"
                predicateClass="System.Web.UI.MobileControls.Adapters.HtmlPageAdapter"
                predicateMethod="DeviceQualifies"
                pageAdapter="System.Web.UI.MobileControls.Adapters.HtmlPageAdapter">

                <control name="System.Web.UI.MobileControls.Panel"             adapter="System.Web.UI.MobileControls.Adapters.HtmlPanelAdapter" />
                <control name="System.Web.UI.MobileControls.Form"              adapter="System.Web.UI.MobileControls.Adapters.HtmlFormAdapter" />
                <control name="System.Web.UI.MobileControls.TextBox"           adapter="System.Web.UI.MobileControls.Adapters.HtmlTextBoxAdapter" />
                <control name="System.Web.UI.MobileControls.Label"             adapter="System.Web.UI.MobileControls.Adapters.HtmlLabelAdapter" />
                <control name="System.Web.UI.MobileControls.LiteralText"       adapter="System.Web.UI.MobileControls.Adapters.HtmlLiteralTextAdapter" />
                <control name="System.Web.UI.MobileControls.Link"              adapter="System.Web.UI.MobileControls.Adapters.HtmlLinkAdapter" />
                <control name="System.Web.UI.MobileControls.Command"           adapter="System.Web.UI.MobileControls.Adapters.HtmlCommandAdapter" />
                <control name="System.Web.UI.MobileControls.PhoneCall"         adapter="System.Web.UI.MobileControls.Adapters.HtmlPhoneCallAdapter" />
                <control name="System.Web.UI.MobileControls.List"              adapter="System.Web.UI.MobileControls.Adapters.HtmlListAdapter" />
                <control name="System.Web.UI.MobileControls.SelectionList"     adapter="System.Web.UI.MobileControls.Adapters.HtmlSelectionListAdapter" />
                <control name="System.Web.UI.MobileControls.ObjectList"        adapter="System.Web.UI.MobileControls.Adapters.HtmlObjectListAdapter" />
                <control name="System.Web.UI.MobileControls.Image"             adapter="System.Web.UI.MobileControls.Adapters.HtmlImageAdapter" />
                <control name="System.Web.UI.MobileControls.BaseValidator"     adapter="System.Web.UI.MobileControls.Adapters.HtmlValidatorAdapter" />
                <control name="System.Web.UI.MobileControls.ValidationSummary" adapter="System.Web.UI.MobileControls.Adapters.HtmlValidationSummaryAdapter" />
                <control name="System.Web.UI.MobileControls.Calendar"          adapter="System.Web.UI.MobileControls.Adapters.HtmlCalendarAdapter" />
                <control name="System.Web.UI.MobileControls.TextView"          adapter="System.Web.UI.MobileControls.Adapters.HtmlTextViewAdapter" />
                <control name="System.Web.UI.MobileControls.MobileControl"     adapter="System.Web.UI.MobileControls.Adapters.HtmlControlAdapter" />
            </device>
            <device name="UpWmlDeviceAdapters"
                inheritsFrom="WmlDeviceAdapters"
                predicateClass="System.Web.UI.MobileControls.Adapters.UpWmlPageAdapter"
                predicateMethod="DeviceQualifies"
                pageAdapter="System.Web.UI.MobileControls.Adapters.UpWmlPageAdapter">
            </device>
            <device name="WmlDeviceAdapters"
                predicateClass="System.Web.UI.MobileControls.Adapters.WmlPageAdapter"
                predicateMethod="DeviceQualifies"
                pageAdapter="System.Web.UI.MobileControls.Adapters.WmlPageAdapter">

                <control name="System.Web.UI.MobileControls.Panel"             adapter="System.Web.UI.MobileControls.Adapters.WmlPanelAdapter" />
                <control name="System.Web.UI.MobileControls.Form"              adapter="System.Web.UI.MobileControls.Adapters.WmlFormAdapter" />
                <control name="System.Web.UI.MobileControls.TextBox"           adapter="System.Web.UI.MobileControls.Adapters.WmlTextBoxAdapter" />
                <control name="System.Web.UI.MobileControls.Label"             adapter="System.Web.UI.MobileControls.Adapters.WmlLabelAdapter" />
                <control name="System.Web.UI.MobileControls.LiteralText"       adapter="System.Web.UI.MobileControls.Adapters.WmlLiteralTextAdapter" />
                <control name="System.Web.UI.MobileControls.Link"              adapter="System.Web.UI.MobileControls.Adapters.WmlLinkAdapter" />
                <control name="System.Web.UI.MobileControls.Command"           adapter="System.Web.UI.MobileControls.Adapters.WmlCommandAdapter" />
                <control name="System.Web.UI.MobileControls.PhoneCall"         adapter="System.Web.UI.MobileControls.Adapters.WmlPhoneCallAdapter" />
                <control name="System.Web.UI.MobileControls.List"              adapter="System.Web.UI.MobileControls.Adapters.WmlListAdapter" />
                <control name="System.Web.UI.MobileControls.SelectionList"     adapter="System.Web.UI.MobileControls.Adapters.WmlSelectionListAdapter" />
                <control name="System.Web.UI.MobileControls.ObjectList"        adapter="System.Web.UI.MobileControls.Adapters.WmlObjectListAdapter" />
                <control name="System.Web.UI.MobileControls.Image"             adapter="System.Web.UI.MobileControls.Adapters.WmlImageAdapter" />
                <control name="System.Web.UI.MobileControls.BaseValidator"     adapter="System.Web.UI.MobileControls.Adapters.WmlValidatorAdapter" />
                <control name="System.Web.UI.MobileControls.ValidationSummary" adapter="System.Web.UI.MobileControls.Adapters.WmlValidationSummaryAdapter" />
                <control name="System.Web.UI.MobileControls.Calendar"          adapter="System.Web.UI.MobileControls.Adapters.WmlCalendarAdapter" />
                <control name="System.Web.UI.MobileControls.TextView"          adapter="System.Web.UI.MobileControls.Adapters.WmlTextViewAdapter" />
                <control name="System.Web.UI.MobileControls.MobileControl"     adapter="System.Web.UI.MobileControls.Adapters.WmlControlAdapter" />
            </device>
            <device name="ChtmlDeviceAdapters"
                inheritsFrom="HtmlDeviceAdapters"
                predicateClass="System.Web.UI.MobileControls.Adapters.ChtmlPageAdapter"
                predicateMethod="DeviceQualifies"
                pageAdapter="System.Web.UI.MobileControls.Adapters.ChtmlPageAdapter">

                <control name="System.Web.UI.MobileControls.Form"              adapter="System.Web.UI.MobileControls.Adapters.ChtmlFormAdapter" />
                <control name="System.Web.UI.MobileControls.Calendar"          adapter="System.Web.UI.MobileControls.Adapters.ChtmlCalendarAdapter" />
                <control name="System.Web.UI.MobileControls.Image"             adapter="System.Web.UI.MobileControls.Adapters.ChtmlImageAdapter" />
                <control name="System.Web.UI.MobileControls.TextBox"           adapter="System.Web.UI.MobileControls.Adapters.ChtmlTextBoxAdapter" />
                <control name="System.Web.UI.MobileControls.SelectionList"     adapter="System.Web.UI.MobileControls.Adapters.ChtmlSelectionListAdapter" />
                <control name="System.Web.UI.MobileControls.Command"           adapter="System.Web.UI.MobileControls.Adapters.ChtmlCommandAdapter" />
                <control name="System.Web.UI.MobileControls.PhoneCall"         adapter="System.Web.UI.MobileControls.Adapters.ChtmlPhoneCallAdapter" />
                <control name="System.Web.UI.MobileControls.Link"              adapter="System.Web.UI.MobileControls.Adapters.ChtmlLinkAdapter" />
            </device>
        </mobileControls>

        <pages>
            <namespaces>
                <add namespace="System" />
                <add namespace="System.Collections" />
                <add namespace="System.Collections.Generic" />
                <add namespace="System.Collections.Specialized" />
                <add namespace="System.ComponentModel.DataAnnotations" />
                <add namespace="System.Configuration" />
                <add namespace="System.Linq" />
                <add namespace="System.Text" />
                <add namespace="System.Text.RegularExpressions" />
                <add namespace="System.Web" />
                <add namespace="System.Web.Caching" />
                <add namespace="System.Web.DynamicData" />
                <add namespace="System.Web.SessionState" />
                <add namespace="System.Web.Security" />
                <add namespace="System.Web.Profile" />
                <add namespace="System.Web.UI" />
                <add namespace="System.Web.UI.WebControls" />
                <add namespace="System.Web.UI.WebControls.WebParts" />
                <add namespace="System.Web.UI.HtmlControls" />
                <add namespace="System.Xml.Linq" />
            </namespaces>

            <controls>
                <add tagPrefix="asp" namespace="System.Web.UI.WebControls.WebParts" assembly="System.Web, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                <add tagPrefix="asp" namespace="System.Web.UI" assembly="System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add tagPrefix="asp" namespace="System.Web.UI.WebControls" assembly="System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add tagPrefix="asp" namespace="System.Web.UI.WebControls.Expressions" assembly="System.Web.Extensions, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add tagPrefix="asp" namespace="System.Web.DynamicData" assembly="System.Web.DynamicData, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
                <add tagPrefix="asp" namespace="System.Web.UI.WebControls" assembly="System.Web.Entity, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089" />
            </controls>
        </pages>

        <protocols>
            <add name="net.tcp" processHandlerType="System.ServiceModel.WasHosting.TcpProcessProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                appDomainHandlerType="System.ServiceModel.WasHosting.TcpAppDomainProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                validate="false" />
            <add name="net.pipe" processHandlerType="System.ServiceModel.WasHosting.NamedPipeProcessProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                appDomainHandlerType="System.ServiceModel.WasHosting.NamedPipeAppDomainProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                validate="false" />
            <add name="net.msmq" processHandlerType="System.ServiceModel.WasHosting.MsmqProcessProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                appDomainHandlerType="System.ServiceModel.WasHosting.MsmqAppDomainProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                validate="false" />
            <add name="msmq.formatname" processHandlerType="System.ServiceModel.WasHosting.MsmqIntegrationProcessProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                appDomainHandlerType="System.ServiceModel.WasHosting.MsmqIntegrationAppDomainProtocolHandler, System.ServiceModel.WasHosting, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089"
                validate="false" />
        </protocols>

        <siteMap>
            <providers>
                <add siteMapFile="web.sitemap" name="AspNetXmlSiteMapProvider"
                    type="System.Web.XmlSiteMapProvider, System.Web, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
            </providers>
        </siteMap>

        <urlMappings enabled="true" />

        <webControls clientScriptsLocation="/aspnet_client/{0}/{1}/" />

        <webParts>
            <personalization>
                <providers>
                    <add connectionStringName="LocalSqlServer"
                        name="AspNetSqlPersonalizationProvider" type="System.Web.UI.WebControls.WebParts.SqlPersonalizationProvider, System.Web, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
                </providers>

                <authorization>
                    <deny users="*" verbs="enterSharedScope" />
                    <allow users="*" verbs="modifyState" />
                </authorization>
            </personalization>

            <transformers>
                <add name="RowToFieldTransformer" type="System.Web.UI.WebControls.WebParts.RowToFieldTransformer" />
                <add name="RowToParametersTransformer" type="System.Web.UI.WebControls.WebParts.RowToParametersTransformer" />
            </transformers>
        </webParts>
				<trace enabled="true" writeToDiagnosticsTrace="true" mostRecent="true" pageOutput="false" />
    </system.web>
    <system.serviceModel>
        <serviceHostingEnvironment>
            <add name="net.tcp" transportConfigurationType="System.ServiceModel.Activation.TcpHostedTransportConfiguration, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
            <add name="net.pipe" transportConfigurationType="System.ServiceModel.Activation.NamedPipeHostedTransportConfiguration, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
            <add name="net.msmq" transportConfigurationType="System.ServiceModel.Activation.MsmqHostedTransportConfiguration, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
            <add name="msmq.formatname" transportConfigurationType="System.ServiceModel.Activation.MsmqIntegrationHostedTransportConfiguration, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" />
        </serviceHostingEnvironment>
    </system.serviceModel>
    <system.xaml.hosting>
        <httpHandlers>
            <add xamlRootElementType="System.ServiceModel.Activities.WorkflowService, System.ServiceModel.Activities, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" httpHandlerType="System.ServiceModel.Activities.Activation.ServiceModelActivitiesActivationHandlerAsync, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
            <add xamlRootElementType="System.Activities.Activity, System.Activities, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35" httpHandlerType="System.ServiceModel.Activities.Activation.ServiceModelActivitiesActivationHandlerAsync, System.ServiceModel.Activation, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35"/>
        </httpHandlers>
    </system.xaml.hosting>

		<system.diagnostics>
  <sharedListeners>
    <add name="PcfLogListener" type="System.Diagnostics.ConsoleTraceListener" />
  </sharedListeners>
  <sources>
    <source name="System.Net">
      <listeners>
        <add name="PcfLogListener"/>
      </listeners>
    </source>
  </sources>
  <trace autoflush="false" indentsize="2">
    <listeners>
      <add name="PcfLogListener"/>
    </listeners>
  </trace>
  <switches>
    <add name="System.Net" value="Information" />
  </switches>
</system.diagnostics>
</configuration>
`
